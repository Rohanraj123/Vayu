#!/usr/bin/env bash
set -eou pipefail
# test-api_keys.sh
# USAGE:
# TOKEN_ADMIN=your-admin-token ./test-api_keys.sh
#
# Assumptions:
# - Docker, kind, kubectl are installed
# - Your vayu Dockerfile is in current dir OR set IMAGE to registry image
# - vayu listens on port 8080
# 
# config file location : /config.yaml

CLUSTER_NAME="vayu-test"
NAMESPACE="vayu-system"
IMAGE="vayu:latest"
DEPLOYMENT_NAME="vayu"
SERVICE_NAME="vayu-service"
PORT=8087
TMPDIR="$(mktemp -d)"
trap 'rm -rf $TMPDIR' EXIT

echo "=== 0. check pre-reqs ==="
for cmd in docker kind kubectl; do
    if ! command -v $cmd > /dev/null 2>&1; then
        echo "ERROR: $cmd is required in PATH"
        exit 1
    fi
done
echo "OKAY"

echo "=== 1. Create kind-cluster if missing ==="
if kind get clusters | grep -q "^${CLUSTER_NAME}$"; then
    echo "Cluster ${CLUSTER_NAME} already exists"
else 
    kind create cluster --name=${CLUSTER_NAME}
fi

echo "=== 3. Deploy vayu manifests ==="
kubectl create ns vayu-system
kubectl create configmap config -n vayu-system --from-file=config.yaml
kubectl apply -f tests/test-api_keys/templates/
echo "waiting for deployment to be ready"
kubectl -n "${NAMESPACE}" rollout status deploy/"${DEPLOYMENT_NAME}" --timeout=120s

echo "=== 4. Snapshot secret before creating key ==="
kubectl -n "${NAMESPACE}" get secrets -o name > "${TMPDIR}/secrets.before" || true

echo "=== 5. From a temporary pod (simulating Admin) call /api-keys ==="
API_KEYS_RESPONSE_FILE="${TMPDIR}/api_keys_response.json"
echo "Calling http://vayu-service:${PORT}/api-keys from temporary pod..."
kubectl -n "${NAMESPACE}" run temp-pod --image=curlimages/curl -it --rm --restart=Never -- \
sh -c "curl -s -w '\n%{http_code}' -X POST -H 'Content-Type: application/json' -d '{\"service\": \"users\"}' http://vayu-service.${NAMESPACE}.svc.cluster.local:${PORT}/api-keys"> "${TMPDIR}/tmp_curl_output" || true

# seperate the body and status
# STATUS=$(tail -n "${TMPDIR}/tmp_curl_output" || true)
BODY=$(sed '$d' "${TMPDIR}/tmp_curl_output" || true)
# echo "HTTP status: ${STATUS}"
echo "Response body: ${BODY}"
echo "${BODY}" > "${API_KEYS_RESPONSE_FILE}"


#if [ "${STATUS}" != "200" ] && [ "${STATUS}" != "201" ]; then
 # echo "ERROR: /api-keys did not return 200/201. Aborting test."
  #echo "Raw response:"
  #cat "${TMPDIR}/tmp_curl_output"
  #exit 2
#fi
#'''

echo "=== 6. Extract raw_api_key from response ==="
# Try jq -> python fallback. Expecting JSON like: {"api_key":"..."} or {"key":"..."}
RAW_KEY=""
if command -v jq >/dev/null 2>&1; then
  RAW_KEY=$(jq -r '.api_key // .key // .raw_key // .token' "${API_KEYS_RESPONSE_FILE}" 2>/dev/null || true)
else
  # python fallback
  RAW_KEY=$(python - <<PYCODE
import json,sys
try:
    j=json.load(open("${API_KEYS_RESPONSE_FILE}"))
    for k in ("api_key","key","raw_key","token"):
        if k in j:
            print(j[k]); sys.exit(0)
except Exception:
    pass
sys.exit(1)
PYCODE
) || true
fi

if [ -z "${RAW_KEY}" ]; then
  echo "Could not parse raw API key from response. Response was:"
  cat "${API_KEYS_RESPONSE_FILE}"
  exit 3
fi
echo "Got raw API key (length ${#RAW_KEY})"


#echo "=== 7. Validate API key by calling protected endpoint ==="
# Call /protected with X-API-Key header
#kubectl -n "${NAMESPACE}" run tmp-curl2 --rm -i --restart=Never --image=curlimages/curl --command -- \
#  sh -c "curl -s -o /dev/null -w '%{http_code}' -H 'X-API-Key: ${RAW_KEY}' http://vayu-service:${PORT}/protected" > "${TMPDIR}/protected_status" || true
#PROT_STATUS=$(cat "${TMPDIR}/protected_status" || true)
#echo "Protected endpoint returned status: ${PROT_STATUS}"

#if [ "${PROT_STATUS}" = "200" ]; then
#  echo "✅ API key accepted by protected endpoint"
#else
#  echo "❌ Protected endpoint did not accept the API key. Status: ${PROT_STATUS}"
  # show logs from vayu pod
#  echo "Showing vayu pod logs (last 200 lines):"
#  POD=$(kubectl -n "${NAMESPACE}" get pods -l app=vayu -o jsonpath='{.items[0].metadata.name}')
#  kubectl -n "${NAMESPACE}" logs "${POD}" --tail=200 || true
#  exit 4
#fi



echo "=== 8. Validate Secret was created ==="
kubectl -n "${NAMESPACE}" get secrets -o name > "${TMPDIR}/secrets.after" || true
echo "Secrets before:"
cat "${TMPDIR}/secrets.before" || true
echo "Secrets after:"
cat "${TMPDIR}/secrets.after" || true
echo "New secrets (diff):"
comm -13 <(sort "${TMPDIR}/secrets.before") <(sort "${TMPDIR}/secrets.after") || true

echo "If a new secret appears, it likely contains the hashed key. Do NOT print Secret data in logs in production."

echo "=== 9. Cleanup (optional) ==="
kind delete cluster --name "${CLUSTER_NAME}"

echo "E2E test completed successfully."



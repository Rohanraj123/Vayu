package apis

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type CreateAPIKeyRequest struct {
	Service string `json:"service"`
}

type CreateAPIKeyResponse struct {
	APIKey string `json:"api_key"`
}

// handler for generating api-keys
func CreateApiKeyHandler(w http.ResponseWriter, r *http.Request, clientset *kubernetes.Clientset) {
	var req CreateAPIKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// generate raw-key
	rawKey := fmt.Sprintf("%x", sha256.Sum256([]byte(req.Service+randomString(8))))

	// hash key for storage
	hash := sha256.Sum256([]byte(rawKey))
	hashedKey := hex.EncodeToString(hash[:])

	// store the haskedKey as a k8s secret
	if err := storeKeyInSecret(req.Service, hashedKey, clientset); err != nil {
		http.Error(w, fmt.Sprintf("failed to store key: %v", err), http.StatusInternalServerError)
		return
	}

	resp := CreateAPIKeyResponse{APIKey: rawKey}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// store api_key into secret of k8s
func storeKeyInSecret(service string, key string, clientset *kubernetes.Clientset) error {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "vayu-api-keys",
			Namespace: "vayu-system",
		},
		StringData: map[string]string{
			service: key,
		},
	}
	_, err := clientset.CoreV1().Secrets("vayu-system").Create(context.TODO(), secret, metav1.CreateOptions{})
	if err != nil {
		// If already exist, update it instead
		_, err = clientset.CoreV1().Secrets("vayu-system").Update(context.TODO(), secret, metav1.UpdateOptions{})
	}

	return err
}

// random string generator
func randomString(bytes int) string {
	characters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789@#$%^&*()")
	b := make([]rune, bytes)
	for i := range b {
		b[i] = characters[i%len(characters)]
	}
	return string(b)
}

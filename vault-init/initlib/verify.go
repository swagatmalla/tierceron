package initlib

import (
	"bitbucket.org/dexterchaney/whoville/validator"
	"bitbucket.org/dexterchaney/whoville/vault-helper/kv"
	"errors"
	"log"
)

// Runs the verification step from data in the seed file
// v is the data contained under the "verification:" tag
// Service name should match credentials in super-secrets
// to verify
// Example
// SpectrumDB:
// 	type: db
// SendGrid:
//	type: SendGridKey
// KeyStore:
// 	type: KeyStore

func verify(mod *kv.Modifier, v map[interface{}]interface{}, logger *log.Logger) ([]string, error) {
	var isValid bool
	var path string
	logger.SetPrefix("[VERIFY]")

	for service, info := range v {
		vType := info.(map[interface{}]interface{})["type"].(string)
		serviceData, err := mod.ReadData("super-secrets/" + service.(string))
		if err != nil {
			return nil, err
		}
		logger.Printf("Verifying %s as type %s\n", service, vType)
		switch vType {
		case "db":
			url := serviceData["url"].(string)
			user := serviceData["user"].(string)
			pass := serviceData["pass"].(string)
			isValid, err = validator.Heartbeat(url, user, pass)
			if err != nil {
				return nil, err
			}
		case "SendGridKey":
			key := serviceData["ApiKey"].(string)
			isValid, err = validator.ValidateSendGrid(key)
			if err != nil {
				return nil, err
			}
		case "KeyStore":
			// path := serviceData["path"].(string)
			// pass := serviceData["pass"].(string)
			isValid = false
		default:
			return nil, errors.New("Invalid verification type: " + vType)
		}

		// Log verification status and write to vault
		logger.Printf("\tverified: %v\n", isValid)
		path = "verification/" + service.(string)
		warn, err := mod.Write(path, map[string]interface{}{
			"type":     vType,
			"verified": isValid,
		})
		if len(warn) > 0 || err != nil {
			return warn, err
		}
	}
	return nil, nil
}
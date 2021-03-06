package publisher

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"time"

	log "github.com/inconshreveable/log15"

	"github.com/netsec-ethz/rains/internal/pkg/keys"
	"github.com/netsec-ethz/rains/internal/pkg/section"
	"github.com/netsec-ethz/rains/internal/pkg/siglib"
	"golang.org/x/crypto/ed25519"
)

//LoadConfig loads configuration information from configPath
func LoadConfig(configPath string) (Config, error) {
	var config Config
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Error("Could not open config file...", "path", configPath, "error", err)
		return Config{}, err
	}
	if err = json.Unmarshal(file, &config); err != nil {
		log.Error("Could not unmarshal json format of config", "error", err)
		return Config{}, err
	}
	config.MetaDataConf.SigSigningInterval *= time.Second
	return config, nil
}

//LoadPrivateKeys reads private keys from the path provided in the config and returns a map from
//PublicKeyID to the corresponding private key data.
func LoadPrivateKeys(path string) (map[keys.PublicKeyID]interface{}, error) {
	var privateKeys []keys.PrivateKey
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Error("Could not open config file...", "path", path, "error", err)
		return nil, err
	}
	if err = json.Unmarshal(file, &privateKeys); err != nil {
		log.Error("Could not unmarshal json format of private keys", "error", err)
		return nil, err
	}
	output := make(map[keys.PublicKeyID]interface{})
	for _, keyData := range privateKeys {
		keyString := keyData.Key.(string)
		privateKey := make([]byte, hex.DecodedLen(len([]byte(keyString))))
		privateKey, err := hex.DecodeString(keyString)
		if err != nil {
			log.Error("Was not able to decode privateKey", "error", err)
			return nil, err
		}
		if len(privateKey) != ed25519.PrivateKeySize {
			log.Error("Private key length is incorrect", "expected", ed25519.PrivateKeySize,
				"actual", len(privateKey))
			return nil, errors.New("incorrect private key length")
		}
		output[keyData.PublicKeyID] = ed25519.PrivateKey(privateKey)
	}
	return output, nil
}

func StorePrivateKey(path string, privateKeys []keys.PrivateKey) error {
	for i, key := range privateKeys {
		privateKeys[i].Key = hex.EncodeToString(key.Key.(ed25519.PrivateKey))
	}
	if encoding, err := json.Marshal(privateKeys); err != nil {
		return err
	} else {
		return ioutil.WriteFile(path, encoding, 0600)
	}
}

//signZone signs the zone and all contained assertions with the zone's private key. It adds the
//subjectZone and context to the contained assertions before signing them and removes them after the
//signatures have been added. It returns an error if it was unable to sign the zone or any of the
//contained assertions.
func signZone(zone *section.Zone, keys map[keys.PublicKeyID]interface{}) error {
	if zone == nil {
		return errors.New("zone is nil")
	}
	zone.DontAddSigInMarshaller()
	if err := signSection(zone, keys); err != nil {
		return err
	}
	zone.AddCtxAndZoneToContent()
	for _, a := range zone.Content {
		if err := signSection(a, keys); err != nil {
			return err
		}
	}
	zone.RemoveCtxAndZoneFromContent()
	zone.AddSigInMarshaller()
	return nil
}

//signShard signs the shard and all contained assertions with the zone's private key. It removes the
//subjectZone and context of the contained assertions after the signatures have been added. It
//returns an error if it was unable to sign the shard or any of the assertions.
func signShard(s *section.Shard, keys map[keys.PublicKeyID]interface{}) error {
	if s == nil {
		return errors.New("shard is nil")
	}
	s.DontAddSigInMarshaller()
	if err := signSection(s, keys); err != nil {
		return err
	}
	s.AddCtxAndZoneToContent()
	for _, a := range s.Content {
		if err := signSection(a, keys); err != nil {
			return err
		}
	}
	s.RemoveCtxAndZoneFromContent()
	s.AddSigInMarshaller()
	return nil
}

//signSection computes the signature data for all contained signatures.
//It returns an error if it was unable to create all signatures on the assertion.
func signSection(s section.WithSigForward, keys map[keys.PublicKeyID]interface{}) error {
	if s == nil {
		return errors.New("section is nil")
	}
	sigs := s.AllSigs()
	s.DeleteAllSigs()
	for _, sig := range sigs {
		if sig.ValidUntil < time.Now().Unix() {
			log.Error("Signature validUntil is in the past")
		} else if ok := siglib.SignSectionUnsafe(s, keys[sig.PublicKeyID], sig); !ok {
			log.Error("Was not able to sign and add the signature", "section", s, "signature", sig)
		} else {
			continue
		}
		return errors.New("Was not able to sign and add the signature")
	}
	return nil
}

package zoneFileParser

import (
	"rains/rainslib"
	"rains/utils/testUtil"
	"testing"

	"fmt"
)

func TestParseZoneFile(t *testing.T) {
	/*test_zone_1 := `
					:Z: . example.com [
						:S: [
							:A: _smtp._tcp [ :srv: mx 25 10 ]
							:A: foobaz [
								:ip4: 192.0.2.33
								:ip6: 2001:db8:cffe:7ea::33
							]
							:A: quuxnorg [
								:ip4: 192.0.3.33
								:ip6: 2001:db8:cffe:7eb::33
							]
						]
					]
					`
	parser := Parser{}
	assertions, err := parser.Decode([]byte(test_zone_1), "")
	if err != nil {
		log.Warn(err.Error())

	} else {
		for _, a := range assertions {
			log.Info(a.Hash())
		}
	}*/
}

func TestEncoder(t *testing.T) {
	nameObjectContent := rainslib.NameObject{
		Name:  "ethz2.ch",
		Types: []rainslib.ObjectType{rainslib.OTIP4Addr, rainslib.OTIP6Addr},
	}
	var ed25519Pkey rainslib.Ed25519PublicKey
	copy(ed25519Pkey[:], []byte("01234567890123456789012345678901"))
	publicKey := rainslib.PublicKey{
		KeySpace: rainslib.RainsKeySpace,
		Type:     rainslib.Ed25519,
		Key:      ed25519Pkey,
	}
	certificate := rainslib.CertificateObject{
		Type:     rainslib.PTTLS,
		HashAlgo: rainslib.Sha256,
		Usage:    rainslib.CUEndEntity,
		Data:     []byte("certData"),
	}
	serviceInfo := rainslib.ServiceInfo{
		Name:     "lookup",
		Port:     49830,
		Priority: 1,
	}

	nameObject := rainslib.Object{Type: rainslib.OTName, Value: nameObjectContent}
	ip6Object := rainslib.Object{Type: rainslib.OTIP6Addr, Value: "2001:0db8:85a3:0000:0000:8a2e:0370:7334"}
	ip4Object := rainslib.Object{Type: rainslib.OTIP4Addr, Value: "127.0.0.1"}
	redirObject := rainslib.Object{Type: rainslib.OTRedirection, Value: "ns.ethz.ch"}
	delegObject := rainslib.Object{Type: rainslib.OTDelegation, Value: publicKey}
	nameSetObject := rainslib.Object{Type: rainslib.OTNameset, Value: rainslib.NamesetExpression("Would be an expression")}
	certObject := rainslib.Object{Type: rainslib.OTCertInfo, Value: certificate}
	serviceInfoObject := rainslib.Object{Type: rainslib.OTServiceInfo, Value: serviceInfo}
	registrarObject := rainslib.Object{Type: rainslib.OTRegistrar, Value: "Registrar information"}
	registrantObject := rainslib.Object{Type: rainslib.OTRegistrant, Value: "Registrant information"}
	infraObject := rainslib.Object{Type: rainslib.OTInfraKey, Value: publicKey}
	extraObject := rainslib.Object{Type: rainslib.OTExtraKey, Value: publicKey}

	signature := rainslib.Signature{
		KeySpace:   rainslib.RainsKeySpace,
		Algorithm:  rainslib.Ed25519,
		ValidSince: 1000,
		ValidUntil: 2000,
		Data:       []byte("SignatureData")}

	containedAssertion := &rainslib.AssertionSection{
		Content: []rainslib.Object{nameObject, ip6Object, ip4Object, redirObject, delegObject, nameSetObject, certObject, serviceInfoObject, registrarObject,
			registrantObject, infraObject, extraObject},
		Context:     "",
		SubjectName: "ethz",
		SubjectZone: "",
		Signatures:  []rainslib.Signature{signature},
	}

	shard := &rainslib.ShardSection{
		Content:     []*rainslib.AssertionSection{containedAssertion},
		Context:     ".",
		SubjectZone: "ch",
		RangeFrom:   "aaa",
		RangeTo:     "zzz",
		Signatures:  []rainslib.Signature{signature},
	}

	zone := &rainslib.ZoneSection{
		Content:     []rainslib.MessageSectionWithSig{containedAssertion, shard},
		Context:     ".",
		SubjectZone: "ch",
		Signatures:  []rainslib.Signature{signature},
	}

	parser := Parser{}
	zoneFile := parser.Encode(zone)
	fmt.Println(zoneFile)
	assertions, err := parser.Decode([]byte(zoneFile), "generatedInTest")
	if err != nil {
		t.Error(err)
	}

	compareAssertion := &rainslib.AssertionSection{
		Content: []rainslib.Object{nameObject, ip6Object, ip4Object, redirObject, delegObject, nameSetObject, certObject, serviceInfoObject, registrarObject,
			registrantObject, infraObject, extraObject},
		Context:     ".",
		SubjectName: "ethz",
		SubjectZone: "ch",
		//no signature is decoded it is generated in rainspub
	}
	testUtil.CheckAssertion(compareAssertion, assertions[0], t)

}

package object

/*func TestNameObjectCompareTo(t *testing.T) {
	nos := sortedNameObjects(9)
	var shuffled []Name
	for _, no := range nos {
		shuffled = append(shuffled, no)
	}
	for i := len(shuffled) - 1; i > 0; i-- {
		j := rand.Intn(i)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}
	sort.Slice(shuffled, func(i, j int) bool { return shuffled[i].CompareTo(shuffled[j]) < 0 })
	for i, no := range nos {
		if !reflect.DeepEqual(no, shuffled[i]) {
			t.Errorf("%d: name objects are in wrong order expected=%v actual%v", i, no, shuffled[i])
		}
	}
}

func TestSignatureAlgorithmTypeString(t *testing.T) {
	var tests = []struct {
		input SignatureAlgorithmType
		want  string
	}{
		{SignatureAlgorithmType(0), "0"},
		{Ed25519, "1"},
		{Ed448, "2"},
		{Ecdsa256, "3"},
		{Ecdsa384, "4"},
	}
	for i, test := range tests {
		if test.input.String() != test.want {
			t.Errorf("%d: Wrong Signature algorithm String value. expected=%v, actual=%v", i, test.want, test.input.String())
		}
	}
}

func TestPublicKeyIDString(t *testing.T) {
	var tests = []struct {
		input PublicKeyID
		want  string
	}{
		{PublicKeyID{}, "AT=0 KS=0 KP=0"},
		{PublicKeyID{Algorithm: Ed25519, KeySpace: RainsKeySpace, KeyPhase: 2}, "AT=1 KS=0 KP=2"},
	}
	for i, test := range tests {
		if test.input.String() != test.want {
			t.Errorf("%d: Wrong public key id String value. expected=%v, actual=%v", i, test.want, test.input.String())
		}
	}
}

func TestPublicKeyString(t *testing.T) {
	var tests = []struct {
		input PublicKey
		want  string
	}{
		{PublicKey{}, "{AT=0 KS=0 KP=0 VS=0 VU=0 data=}"},
		{
			PublicKey{
				PublicKeyID: PublicKeyID{
					Algorithm: Ed25519,
					KeySpace:  RainsKeySpace,
					KeyPhase:  1,
				},
				ValidSince: 1,
				ValidUntil: 2,
				Key:        ed25519.PublicKey([]byte("PublicKeyData"))},
			"{AT=1 KS=0 KP=1 VS=1 VU=2 data=5075626c69634b657944617461}",
		},
	}
	for i, test := range tests {
		if test.input.String() != test.want {
			t.Errorf("%d: Wrong public key String value. expected=%v, actual=%v", i, test.want, test.input.String())
		}
	}
}

func TestCertTypeString(t *testing.T) {
	var tests = []struct {
		input Certificate
		want  string
	}{
		{Certificate{}, "{0 0 0 }"},
		{Certificate{Type: PTTLS, Usage: CUEndEntity, HashAlgo: Sha256, Data: []byte("testCert")}, "{1 3 1 7465737443657274}"},
	}
	for i, test := range tests {
		if test.input.String() != test.want {
			t.Errorf("%d: Wrong cert type String value. expected=%v, actual=%v", i, test.want, test.input.String())
		}
	}
}

func TestPublicKeyIDHash(t *testing.T) {
	var tests = []struct {
		input PublicKeyID
		want  string
	}{
		{PublicKeyID{}, "0,0,0"},
		{PublicKeyID{Algorithm: Ed25519, KeySpace: RainsKeySpace, KeyPhase: 2}, "1,0,2"},
	}
	for i, test := range tests {
		if test.input.Hash() != test.want {
			t.Errorf("%d: Wrong Public key id Hash value. expected=%v, actual=%v", i, test.want, test.input.Hash())
		}
	}
}

func TestPublicKeyHash(t *testing.T) {
	var tests = []struct {
		input PublicKey
		want  string
	}{
		{PublicKey{}, "0,0,0,0,0,"},
		{
			PublicKey{
				PublicKeyID: PublicKeyID{
					Algorithm: Ed25519,
					KeySpace:  RainsKeySpace,
					KeyPhase:  1,
				},
				ValidSince: 1,
				ValidUntil: 2,
				Key:        ed25519.PublicKey([]byte("PublicKeyData"))},
			"1,0,1,1,2,5075626c69634b657944617461",
		},
	}
	for i, test := range tests {
		if test.input.Hash() != test.want {
			t.Errorf("%d: Wrong public key String value. expected=%v, actual=%v", i, test.want, test.input.Hash())
		}
	}
}

func TestPublicKeyCompareTo(t *testing.T) {
	pks := sortedPublicKeys(9)
	var shuffled []PublicKey
	for _, pk := range pks {
		shuffled = append(shuffled, pk)
	}
	for i := len(shuffled) - 1; i > 0; i-- {
		j := rand.Intn(i)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}
	sort.Slice(shuffled, func(i, j int) bool { return shuffled[i].CompareTo(shuffled[j]) < 0 })
	for i, pk := range pks {
		if !reflect.DeepEqual(pk, shuffled[i]) {
			t.Errorf("%d: name objects are in wrong order expected=%v actual%v", i, pk, shuffled[i])
		}
	}
	pk1 := pks[0]
	pk1.Key = []byte{}
	if pk1.CompareTo(pks[0]) != 0 {
		t.Error("Error case was not hit")
	}
	if pks[0].CompareTo(pk1) != 0 {
		t.Error("Error case was not hit")
	}

	pk1.KeySpace = KeySpaceID(1)
	if pk1.CompareTo(pks[0]) != 1 {
		t.Error("key space comparison")
	}
	pk1.KeySpace = KeySpaceID(-1)
	if pk1.CompareTo(pks[0]) != -1 {
		t.Error("key space comparison")
	}
}

func TestCertificateCompareTo(t *testing.T) {
	certs := sortedCertificates(9)
	var shuffled []Certificate
	for _, cert := range certs {
		shuffled = append(shuffled, cert)
	}
	for i := len(shuffled) - 1; i > 0; i-- {
		j := rand.Intn(i)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}
	sort.Slice(shuffled, func(i, j int) bool { return shuffled[i].CompareTo(shuffled[j]) < 0 })
	for i, cert := range certs {
		if !reflect.DeepEqual(cert, shuffled[i]) {
			t.Errorf("%d: name objects are in wrong order expected=%v actual%v", i, cert, shuffled[i])
		}
	}
}

func TestServiceInfoCompareTo(t *testing.T) {
	sis := sortedServiceInfo(5)
	var shuffled []ServiceInfo
	for _, si := range sis {
		shuffled = append(shuffled, si)
	}
	for i := len(shuffled) - 1; i > 0; i-- {
		j := rand.Intn(i)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}
	sort.Slice(shuffled, func(i, j int) bool { return shuffled[i].CompareTo(shuffled[j]) < 0 })
	for i, si := range sis {
		if !reflect.DeepEqual(si, shuffled[i]) {
			t.Errorf("%d: name objects are in wrong order expected=%v actual%v", i, si, shuffled[i])
		}
	}
}

func TestObjectCompareTo(t *testing.T) {
	objs := sortedObjects(13)
	var shuffled []Object
	for _, obj := range objs {
		shuffled = append(shuffled, obj)
	}
	for i := len(shuffled) - 1; i > 0; i-- {
		j := rand.Intn(i)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}
	sort.Slice(shuffled, func(i, j int) bool { return shuffled[i].CompareTo(shuffled[j]) < 0 })
	for i, obj := range objs {
		if !reflect.DeepEqual(obj, shuffled[i]) {
			t.Errorf("%d: name objects are in wrong order expected=%v actual%v", i, obj, shuffled[i])
		}
	}
	//Test error cases
	obj1 := objs[0]
	obj1.Value = ""
	if obj1.CompareTo(objs[0]) != 0 {
		t.Error("Error case was not hit")
	}
	obj1.Value = PublicKey{}
	if obj1.CompareTo(objs[0]) != 0 {
		t.Error("Error case was not hit")
	}
	obj1.Value = Certificate{}
	if obj1.CompareTo(objs[0]) != 0 {
		t.Error("Error case was not hit")
	}
	obj1.Value = ServiceInfo{}
	if obj1.CompareTo(objs[0]) != 0 {
		t.Error("Error case was not hit")
	}
	obj1.Value = NamesetExpr("Test")
	if obj1.CompareTo(objs[0]) != 0 {
		t.Error("Error case was not hit")
	}
	obj1.Value = 5
	if obj1.CompareTo(objs[0]) != 0 {
		t.Error("Error case was not hit")
	}
	if objs[0].CompareTo(obj1) != 0 {
		t.Error("Error case was not hit")
	}
}
func TestObjectString(t *testing.T) {
	obj := GetAllValidObjects()
	var tests = []struct {
		input Object
		want  string
	}{
		{Object{}, "OT:0 OV:<nil>"},
		{obj[0], "OT:1 OV:{example.com [3 2]}"},
		{obj[1], "OT:2 OV:2001:db8::"},
		{obj[2], "OT:3 OV:192.0.2.0"},
		{obj[3], "OT:4 OV:example.com"},
		{obj[4], fmt.Sprintf("OT:5 OV:%s", obj[4].Value.(PublicKey).String())},
		{obj[5], "OT:6 OV:Would be an expression"},
		{obj[6], fmt.Sprintf("OT:7 OV:%s", obj[6].Value.(Certificate).String())},
		{obj[7], "OT:8 OV:{srvName 49830 1}"},
		{obj[8], "OT:9 OV:Registrar information"},
		{obj[9], "OT:10 OV:Registrant information"},
		{obj[10], fmt.Sprintf("OT:11 OV:%s", obj[10].Value.(PublicKey).String())},
		{obj[11], fmt.Sprintf("OT:12 OV:%s", obj[11].Value.(PublicKey).String())},
		{obj[12], fmt.Sprintf("OT:13 OV:%s", obj[12].Value.(PublicKey).String())},
	}
	for i, test := range tests {
		if test.input.String() != test.want {
			t.Errorf("%d: Wrong Object String(). expected=%v, actual=%v", i, test.want, test.input.String())
		}
	}
}

func TestObjectSort(t *testing.T) {
	objTypes := []Type{OTNextKey, OTExtraKey, OTInfraKey, OTRegistrant, OTRegistrar, OTServiceInfo, OTCertInfo, OTNameset, OTDelegation, OTRedirection,
		OTIP4Addr, OTIP6Addr, OTName}
	expected := []Type{OTName, OTIP6Addr, OTIP4Addr, OTRedirection, OTDelegation, OTNameset, OTCertInfo, OTServiceInfo, OTRegistrar, OTRegistrant,
		OTInfraKey, OTExtraKey, OTNextKey}
	obj := Object{Type: OTName, Value: Name{Name: "", Types: objTypes}}
	expectedObj := Object{Type: OTName, Value: Name{Name: "", Types: expected}}
	obj.Sort()
	if !reflect.DeepEqual(obj, expectedObj) {
		t.Errorf("name objects are in wrong order after obj.Sort() expected=%v actual%v", expectedObj, obj)
	}
	//error case
	obj = Object{Type: OTExtraKey, Value: ""}
	obj.Sort()
}*/

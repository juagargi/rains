@0xfb2d77234707241e; # unique file ID, generated by `capnp id`
using Go = import "/go.capnp";
$Go.package("proto");
$Go.import("rains/proto");

struct RainsMessage  {
    #RainsMessage contains the data of a message
	token           @0 :Data;
	content         @1 :List(MessageSection);
	signatures      @2 :List(Signature);
	capabilities    @3 :List(Text);
}

struct MessageSection  {
    union {
        assertion           @0 :AssertionSection;     
        shard               @1 :ShardSection; 
        zone                @2 :ZoneSection;      
        query               @3 :QuerySection;      
        notification        @4 :NotificationSection;
        addressQuery        @5 :AddressQuerySection;
        addressAssertion    @6 :AddressAssertionSection;
    }         
}

const noCapability :Text = "";
const tLSOverTCP :Text = "urn:x-rains:tlssrv";


interface MessageSectionWithSig extends(Interval) {
    #MessageSectionWithSig can be either an Assertion, Shard or Zone
	sigs            @0 () -> (sig :Signature);
	addSig          @1 (sig :Signature);
	deleteSig       @2 (int :Int32);
	deleteAllSigs   @3 ();
   #TODO CFE what is the syntax of a method without arguments
	getContext      @4 () -> (context :Text);
	getSubjectZone  @5 () -> (zone :Text);
	createStub      @6 () -> (section :MessageSectionWithSig);
	validFrom       @7 () -> (validFrom :Int64);
	validUntil      @8 () -> (validUntil :Int64);
	hash            @9 () -> (hash :Text);
}

interface Interval  {
    #Interval defines an interval over strings
	begin @0 () -> (begin :Text);
    end     @1 () -> (end :Text);
}

interface Hashable  {
    #Hashable can be implemented by objects that are not natively hashable.
	hash @0 () -> (hash :Text);
}

struct AssertionSection  {
    #AssertionSection contains information about the assertion
	subjectName @0 :Text;
	content     @1 :List(Obj);
	signatures  @2 :List(Signature);
	subjectZone @3 :Text;
	context     @4 :Text;
}


struct ShardSection  {
    #ShardSection contains information about the shard
	content     @0 :List(AssertionSection);
	signatures  @1 :List(Signature);
	subjectZone @2 :Text;
	context     @3 :Text;
	rangeFrom   @4 :Text;
	rangeTo     @5 :Text;
}


struct ZoneSection  {
    #ZoneSection contains information about the zone
	signatures  @0  :List(Signature);
	subjectZone @1  :Text;
	context     @2  :Text;
	content     @3  :List(MessageSectionWithSig);
}

struct QuerySection  {
    #QuerySection contains information about the query
	token   @0      :Data;
	name    @1      :Text;
	context @2      :Text;
	type    @3      :ObjectType;
	expires @4      :Int64; #time when this query expires represented as the number of seconds elapsed since January 1, 1970 UTC
	options @5      :List(QueryOption);
}

enum QueryOption {
	minE2ELatency            @0;
	minLastHopAnswerSize     @1;
	minInfoLeakage           @2;
	cachedAnswersOnly        @3;
	expiredAssertionsOk      @4;
	tokenTracing             @5;
	noVerificationDelegation @6;
	noProactiveCaching       @7;
}

enum ObjectType {
	oTName        @0;
	oTIP6Addr     @1;
	oTIP4Addr     @2;
	oTRedirection @3;
	oTDelegation  @4;
	oTNameset     @5;
	oTCertInfo    @6;
	oTServiceInfo @7;
	oTRegistrar   @8;
	oTRegistrant  @9;
	oTInfraKey    @10;
	oTExtraKey    @11;
}

struct SubjectAddr  {
	addressFamily @0    :Text;
	prefixLength  @1    :UInt32;
	address       @2    :Text;
}

struct AddressAssertionSection  {
    #AddressAssertionSection contains information about the address assertion
	subjectAddr @0  :SubjectAddr;
	content     @1  :List(Obj);
	signatures  @2  :List(Signature);
	context     @3  :Text;
}


struct AddressZoneSection  {
    #AddressZoneSection contains information about the address zone
	subjectAddr @0  :SubjectAddr;
	signatures  @1  :List(Signature);
	context     @2  :Text;
	content     @3  :List(AddressAssertionSection);
}

struct AddressQuerySection  {
    #AddressQuerySection contains information about the address query
	subjectAddr @0 :SubjectAddr;
	token       @1 :Data;
	context     @2 :Text;
	types       @3 :List(Int32);
	expires     @4 :Int64;
	options     @5 :List(QueryOption);
}


struct NotificationSection  {
    #NotificationSection contains information about the notification
	token @0    :Data;
	type  @1    :NotificationType;
	data  @2    :Text;
}

enum NotificationType {
	heartbeat          @0;
	capHashNotKnown    @1;
	badMessage         @2;
	rcvInconsistentMsg @3;
	noAssertionsExist  @4;
	msgTooLarge        @5;
	unspecServerErr    @6;
	serverNotCapable   @7;
	noAssertionAvail   @8;
}


struct Signature  {
    #Signature on a Rains message or section
	keySpace   @0 :KeySpaceID;
	algorithm  @1 :SignatureAlgorithmType;
	validSince @2 :Int64;
	validUntil @3 :Int64;
	data       @4 :Data;
}

enum KeySpaceID {
#KeySpaceID identifies a key space
	rainsKeySpace @0;
}

enum SignatureAlgorithmType {
#SignatureAlgorithmType specifies a signature algorithm type
	ed25519  @0;
	ed448    @1;
	ecdsa256 @2;
	ecdsa384 @3;
}

enum HashAlgorithmType {
#HashAlgorithmType specifies a hash algorithm type
	noHashAlgo @0;
	sha256     @1;
	sha384     @2;
	sha512     @3;
}

struct PublicKey  {
    #PublicKey contains information about a public key
	type       @0 :SignatureAlgorithmType;
	key        @1 :Data;
	validFrom  @2 :Int64;
	validUntil @3 :Int64;
}

struct CertificateObject  {
    #CertificateObject contains certificate information
	type     @0 :ProtocolType;
	usage    @1 :CertificateUsage;
	hashAlgo @2 :HashAlgorithmType;
	data     @3 :Data;
}

enum ProtocolType {
	pTUnspecified @0;
	pTTLS         @1;
}

enum CertificateUsage {
	cUTrustAnchor @0;
	cUEndEntity   @1;
}

struct ServiceInfo  {
    #ServiceInfo contains information how to access a named service
	name     @0 :Text;
	port     @1 :UInt16;
	priority @2 :UInt32;
}

struct Obj  {
    #Object is a container for different values determined by the given type.
	type  @0    :ObjectType;
	value   :union {
        ip4 @1 :Text;
        ip6 @2 :Text;
    }
}



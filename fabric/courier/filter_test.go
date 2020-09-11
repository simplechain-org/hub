package courier

import (
	"encoding/hex"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/common"
)

func TestGetPrepareCrossTxs(t *testing.T) {
	example := "0a4608041220c316f12ef6e89df8be3faa0bd501379f2758d665c5f5a990e2bd5ed15cc3ae4b1a20c8a81ed81949c2031b8120163a424f43d4ac991b92e0aefd09d031ac2e02c89412921f0a8f1f0ac41e0ac2070a6608031a0b0882a19cfa0510c8ee914222096d796368616e6e656c2a40306636336233623261363334633533633338323965613836353831393962313930656431646236623730383837353431356437323063343438663030323265663a08120612046d79636312d7060aba060a074f7267324d535012ae062d2d2d2d2d424547494e2043455254494649434154452d2d2d2d2d0a4d4949434b7a4343416447674177494241674952414b6f6e694c5071716767575365777677352b7946783477436759494b6f5a497a6a304541774977637a454c0a4d416b474131554542684d4356564d78457a415242674e5642416754436b4e6862476c6d62334a7561574578466a415542674e564241635444564e68626942470a636d467559326c7a593238784754415842674e5642416f54454739795a7a49755a586868625842735a53356a623230784844416142674e5642414d5445324e680a4c6d39795a7a49755a586868625842735a53356a623230774868634e4d6a41774f4449334d4445304d5441775768634e4d7a41774f4449314d4445304d5441770a576a42734d517377435159445651514745774a56557a45544d4245474131554543424d4b5132467361575a76636d3570595445574d4251474131554542784d4e0a5532467549455a795957356a61584e6a627a45504d4130474131554543784d47593278705a5735304d523877485159445651514444425a425a473170626b42760a636d63794c6d56345957317762475575593239744d466b77457759484b6f5a497a6a3043415159494b6f5a497a6a304441516344516741457152386e2b5752640a656c77544a5348374b6a6f42554d36302f7a49593179334c57505a785563742f6a784172336954657354423935454a66435271494a38515437315579767330300a71336d576b7342736577446378614e4e4d45737744675944565230504151482f42415144416765414d41774741315564457745422f7751434d4141774b7759440a5652306a42435177496f41672b4e4c4f31716c4e3256396462396a5a4a6a7a38786c42366678736167324b6a316d376f652f547069693877436759494b6f5a490a7a6a3045417749445341417752514968414a2f724868706858736e755549774557627171514935326a4f34304446317a425151494348534443524a52416942360a676b71583979426e385a7a687272454147776b34685462746c38705675506d6953644a2f437a527375773d3d0a2d2d2d2d2d454e442043455254494649434154452d2d2d2d2d0a1218ffbfc9c733aca7d3c16c53b3c2dee5920abda8b1d2662adf12fc160af9160ad7060aba060a074f7267324d535012ae062d2d2d2d2d424547494e2043455254494649434154452d2d2d2d2d0a4d4949434b7a4343416447674177494241674952414b6f6e694c5071716767575365777677352b7946783477436759494b6f5a497a6a304541774977637a454c0a4d416b474131554542684d4356564d78457a415242674e5642416754436b4e6862476c6d62334a7561574578466a415542674e564241635444564e68626942470a636d467559326c7a593238784754415842674e5642416f54454739795a7a49755a586868625842735a53356a623230784844416142674e5642414d5445324e680a4c6d39795a7a49755a586868625842735a53356a623230774868634e4d6a41774f4449334d4445304d5441775768634e4d7a41774f4449314d4445304d5441770a576a42734d517377435159445651514745774a56557a45544d4245474131554543424d4b5132467361575a76636d3570595445574d4251474131554542784d4e0a5532467549455a795957356a61584e6a627a45504d4130474131554543784d47593278705a5735304d523877485159445651514444425a425a473170626b42760a636d63794c6d56345957317762475575593239744d466b77457759484b6f5a497a6a3043415159494b6f5a497a6a304441516344516741457152386e2b5752640a656c77544a5348374b6a6f42554d36302f7a49593179334c57505a785563742f6a784172336954657354423935454a66435271494a38515437315579767330300a71336d576b7342736577446378614e4e4d45737744675944565230504151482f42415144416765414d41774741315564457745422f7751434d4141774b7759440a5652306a42435177496f41672b4e4c4f31716c4e3256396462396a5a4a6a7a38786c42366678736167324b6a316d376f652f547069693877436759494b6f5a490a7a6a3045417749445341417752514968414a2f724868706858736e755549774557627171514935326a4f34304446317a425151494348534443524a52416942360a676b71583979426e385a7a687272454147776b34685462746c38705675506d6953644a2f437a527375773d3d0a2d2d2d2d2d454e442043455254494649434154452d2d2d2d2d0a1218ffbfc9c733aca7d3c16c53b3c2dee5920abda8b1d2662adf129c100a220a200a1e0801120612046d7963631a120a06696e766f6b650a01610a01620a02313012f50f0ae9010a200d373bb0b7b5d6f929e0b519058fa3435517c79a35a44d346500d2eec26d9e8312c4010a4512140a046c736363120c0a0a0a046d79636312020803122d0a046d79636312250a070a0161120208030a070a0162120208031a070a01611a0239301a080a01621a0332313012690a046d7963631240306636336233623261363334633533633338323965613836353831393962313930656431646236623730383837353431356437323063343438663030323265661a0b6576745472616e73666572221261207472616e7366657220746f20622031301a0308c801220b12046d7963631a03312e301281070ab6060a074f7267314d535012aa062d2d2d2d2d424547494e2043455254494649434154452d2d2d2d2d0a4d4949434b44434341632b674177494241674952414c4179424434675a48564176755146312b6c5035513477436759494b6f5a497a6a304541774977637a454c0a4d416b474131554542684d4356564d78457a415242674e5642416754436b4e6862476c6d62334a7561574578466a415542674e564241635444564e68626942470a636d467559326c7a593238784754415842674e5642416f54454739795a7a45755a586868625842735a53356a623230784844416142674e5642414d5445324e680a4c6d39795a7a45755a586868625842735a53356a623230774868634e4d6a41774f4449334d4445304d5441775768634e4d7a41774f4449314d4445304d5441770a576a42714d517377435159445651514745774a56557a45544d4245474131554543424d4b5132467361575a76636d3570595445574d4251474131554542784d4e0a5532467549455a795957356a61584e6a627a454e4d4173474131554543784d456347566c636a45664d4230474131554541784d576347566c636a417562334a6e0a4d53356c654746746347786c4c6d4e766254425a4d424d4742797147534d34394167454743437147534d34394177454841304941424f57695773436e447344370a756d4e653433357143412f3052725637306b674f794a2b694c5970524863487547517338714175544e7a5637652b2b484f3665465448306346456d742f6376340a2f644948784731365361696a5454424c4d41344741315564447745422f775145417749486744414d42674e5648524d4241663845416a41414d437347413155640a4977516b4d434b41494c555241795959506739596c544c4a5a4756724f4566696c446d574b74483959475461773768436b2f4e494d416f4743437147534d34390a42414d43413063414d45514349433550344262776b766861314e513372526a734d4a46594c646141364a665a374579746b785a75494e6f37416941437962376d0a4766486c444243676a3561316f30662b33717a577149646a474850674b426f4679762f5936673d3d0a2d2d2d2d2d454e442043455254494649434154452d2d2d2d2d0a1246304402204aea84f7a5d87fa07a8861c203750f6f87de0510fd5a9f1c4549cbc07ba32a580220656dda4dff52b3b792540a735b18c6cbaa8fe6266837a9eb07a62f1ebd6a6e581282070ab6060a074f7267324d535012aa062d2d2d2d2d424547494e2043455254494649434154452d2d2d2d2d0a4d4949434b54434341632b674177494241674952414e6d41796152584f684d572b7467436e69477578675177436759494b6f5a497a6a304541774977637a454c0a4d416b474131554542684d4356564d78457a415242674e5642416754436b4e6862476c6d62334a7561574578466a415542674e564241635444564e68626942470a636d467559326c7a593238784754415842674e5642416f54454739795a7a49755a586868625842735a53356a623230784844416142674e5642414d5445324e680a4c6d39795a7a49755a586868625842735a53356a623230774868634e4d6a41774f4449334d4445304d5441775768634e4d7a41774f4449314d4445304d5441770a576a42714d517377435159445651514745774a56557a45544d4245474131554543424d4b5132467361575a76636d3570595445574d4251474131554542784d4e0a5532467549455a795957356a61584e6a627a454e4d4173474131554543784d456347566c636a45664d4230474131554541784d576347566c636a417562334a6e0a4d69356c654746746347786c4c6d4e766254425a4d424d4742797147534d34394167454743437147534d3439417745484130494142486b396a466b5a696c69590a4467793931466f74632f446c30454c2b775a5236785144795851646b34336f5930336c776d5a3776656e6c636b4741426663445a394669436244426a345738770a64576e674d4b617961622b6a5454424c4d41344741315564447745422f775145417749486744414d42674e5648524d4241663845416a41414d437347413155640a4977516b4d434b4149506a537a74617054646c6658572f59325359382f4d5a51656e3862476f4e696f395a753648763036596f764d416f4743437147534d34390a42414d43413067414d45554349514434426262775262474c704a6e434233316b577a4e525678366878356d4a7579502b6f52537736744d7870414967632b50490a746e5068375a52764256314e4a6d4b64764843656f787254734971396e593376333049764d49553d0a2d2d2d2d2d454e442043455254494649434154452d2d2d2d2d0a12473045022100a406f5a7eecbf4b8909ac5b65c835b00133f4e96a662fbfa758cb43fdef87ee502206f7f5caa38758fbf10302cc6be296dca3e7a1431bd7abc2c49421210807f05401246304402203c74752e4f6a389d319ba4a0d022c58b64a07d2b4fdaecec10dd4c2e02764d68022066b83af9402cec194f655d2f9ad5c797c80492d84d8231c12d0ee2111734dcd41ab3070a81070a040a02080212f8060aad060a90060a0a4f7264657265724d53501281062d2d2d2d2d424547494e2043455254494649434154452d2d2d2d2d0a4d4949434444434341624f674177494241674952414e42566435725a51784268546c6366463570416f4a3877436759494b6f5a497a6a3045417749776154454c0a4d416b474131554542684d4356564d78457a415242674e5642416754436b4e6862476c6d62334a7561574578466a415542674e564241635444564e68626942470a636d467559326c7a593238784644415342674e5642416f54433256345957317762475575593239744d52637746515944565151444577356a5953356c654746740a6347786c4c6d4e7662544165467730794d4441344d6a63774d5451784d4442614677307a4d4441344d6a55774d5451784d4442614d466778437a414a42674e560a42415954416c56544d524d77455159445651514945777044595778705a6d3979626d6c684d52597746415944565151484577315459573467526e4a68626d4e700a63324e764d527777476759445651514445784e76636d526c636d56794c6d56345957317762475575593239744d466b77457759484b6f5a497a6a3043415159490a4b6f5a497a6a304441516344516741454b32685737597535465335737459755a64617561425770656e6650496b522f505957717351316439416a77425a314e300a58334c693173316237736f7544735263447a337a396477566c57686b44736778744e6b4a52614e4e4d45737744675944565230504151482f42415144416765410a4d41774741315564457745422f7751434d4141774b7759445652306a42435177496f41674d30506579645a7a6d4e3569573669776d46526974754c58394972640a463258634b38412f613457497a724577436759494b6f5a497a6a304541774944527741775241496750596c75766f5848764b7657577147694d7774495253712b0a634430634237683036347143594f32365254454349437a2b6e62556568726d6a4d4f672f5a41714e694f772b6a33676b51695033355a7066382b4f4f483636620a2d2d2d2d2d454e442043455254494649434154452d2d2d2d2d0a12188f04b960c139d17d5d3a57c51fda99e0cd2802ecfa5cd4331246304402203d16e65f1902c3b3e982c3e363918d453eea6886e25b3b0e7bf50b29cbde63cc02203e5c4be740bcc8870f2613c4bbea9ef371e8ec1722e1c222cf9090808b7d616f0a040a0208020a01000a000a220a2038fb5437e6af9387bec51324df8f9aaae0524e507089fa3e31d079c1cf7a7d32"

	rawBlock, err := hex.DecodeString(example)
	if err != nil {
		t.Fatal(err)
	}

	var block common.Block
	if err := proto.Unmarshal(rawBlock, &block); err != nil {
		t.Fatal(err)
	}

	preTxs, err := GetPrepareCrossTxs(&block, func(ev string) bool {
		return true
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(preTxs) != 1 {
		t.Fatalf("len(preTxs), want: 1, got: %d", len(preTxs))
	}

	if preTxs[0].BlockNumber != 4 {
		t.Fatalf("Number, want: 4, got: %d", preTxs[0].BlockNumber)
	}

	if preTxs[0].TxID != "0f63b3b2a634c53c3829ea8658199b190ed1db6b708875415d720c448f0022ef" {
		t.Fatalf("TxID, want: <0f63b3b2a634c53c3829ea8658199b190ed1db6b708875415d720c448f0022ef>, got: <%s>", preTxs[0].TxID)
	}

	if preTxs[0].TimeStamp.Seconds != 1598492802 &&
		preTxs[0].TimeStamp.Nanos != 138704712 {
		t.Fatalf("TimeStamp, want: <seconds=1598492802 nanos=138704712>, "+
			"got: <seconds=%d nanos=%d>", preTxs[0].TimeStamp.Seconds, preTxs[0].TimeStamp.Nanos)
	}

	if preTxs[0].EventName != "evtTransfer" {
		t.Fatalf("EventName, want: <evtTransfer>, got: <%s>", preTxs[0].EventName)
	}

	if string(preTxs[0].Payload) != "a transfer to b 10" {
		t.Fatalf("Payload, want: <evtTransfer>, got: <%s>", string(preTxs[0].Payload))
	}
}

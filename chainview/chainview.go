package chainview

import (
	"bytes"
	"context"
	"fmt"
	"github.com/simplechain-org/crosshub/core"
	"github.com/simplechain-org/crosshub/repo"
	"github.com/simplechain-org/go-simplechain"
	"github.com/simplechain-org/go-simplechain/crypto"
	"github.com/simplechain-org/go-simplechain/crypto/ecdsa"
	"github.com/simplechain-org/go-simplechain/log"
	"github.com/simplechain-org/go-simplechain/rpc"
	"time"

	"github.com/simplechain-org/go-simplechain/accounts/abi"
	"github.com/simplechain-org/go-simplechain/core/types"
	"math/big"

	"github.com/simplechain-org/go-simplechain/common"
	"github.com/simplechain-org/go-simplechain/common/hexutil"
	"github.com/simplechain-org/go-simplechain/ethclient"
)

var abiParsed abi.ABI
var CrossAbi = "0x5b0a097b0a090922696e70757473223a205b0a0909097b0a0909090922696e7465726e616c54797065223a202275696e7438222c0a09090909226e616d65223a2022707572706f7365222c0a090909092274797065223a202275696e7438220a0909097d2c0a0909097b0a0909090922696e7465726e616c54797065223a2022616464726573732070617961626c65222c0a09090909226e616d65223a2022616e63686f72222c0a090909092274797065223a202261646472657373220a0909097d2c0a0909097b0a0909090922696e7465726e616c54797065223a202275696e74323536222c0a09090909226e616d65223a2022726577617264222c0a090909092274797065223a202275696e74323536220a0909097d0a09095d2c0a0909226e616d65223a2022616363756d756c61746552657761726473222c0a0909226f757470757473223a205b5d2c0a09092273746174654d75746162696c697479223a20226e6f6e70617961626c65222c0a09092274797065223a202266756e6374696f6e220a097d2c0a097b0a090922696e70757473223a205b5d2c0a09092273746174654d75746162696c697479223a20226e6f6e70617961626c65222c0a09092274797065223a2022636f6e7374727563746f72220a097d2c0a097b0a090922616e6f6e796d6f7573223a2066616c73652c0a090922696e70757473223a205b0a0909097b0a0909090922696e6465786564223a2066616c73652c0a0909090922696e7465726e616c54797065223a202275696e7438222c0a09090909226e616d65223a2022707572706f7365222c0a090909092274797065223a202275696e7438220a0909097d2c0a0909097b0a0909090922696e6465786564223a2066616c73652c0a0909090922696e7465726e616c54797065223a202261646472657373222c0a09090909226e616d65223a2022616e63686f72222c0a090909092274797065223a202261646472657373220a0909097d2c0a0909097b0a0909090922696e6465786564223a2066616c73652c0a0909090922696e7465726e616c54797065223a202275696e74323536222c0a09090909226e616d65223a2022726577617264222c0a090909092274797065223a202275696e74323536220a0909097d0a09095d2c0a0909226e616d65223a2022416363756d756c61746552657761726473222c0a09092274797065223a20226576656e74220a097d2c0a097b0a090922696e70757473223a205b0a0909097b0a0909090922696e7465726e616c54797065223a202275696e7438222c0a09090909226e616d65223a2022707572706f7365222c0a090909092274797065223a202275696e7438220a0909097d2c0a0909097b0a0909090922696e7465726e616c54797065223a2022616464726573735b5d222c0a09090909226e616d65223a20225f616e63686f7273222c0a090909092274797065223a2022616464726573735b5d220a0909097d0a09095d2c0a0909226e616d65223a2022616464416e63686f7273222c0a0909226f757470757473223a205b5d2c0a09092273746174654d75746162696c697479223a20226e6f6e70617961626c65222c0a09092274797065223a202266756e6374696f6e220a097d2c0a097b0a090922616e6f6e796d6f7573223a2066616c73652c0a090922696e70757473223a205b0a0909097b0a0909090922696e6465786564223a2066616c73652c0a0909090922696e7465726e616c54797065223a202275696e7438222c0a09090909226e616d65223a2022707572706f7365222c0a090909092274797065223a202275696e7438220a0909097d0a09095d2c0a0909226e616d65223a2022416464416e63686f7273222c0a09092274797065223a20226576656e74220a097d2c0a097b0a090922696e70757473223a205b0a0909097b0a0909090922696e7465726e616c54797065223a202275696e7438222c0a09090909226e616d65223a2022707572706f7365222c0a090909092274797065223a202275696e7438220a0909097d2c0a0909097b0a0909090922696e7465726e616c54797065223a202275696e74323536222c0a09090909226e616d65223a20226d617856616c7565222c0a090909092274797065223a202275696e74323536220a0909097d2c0a0909097b0a0909090922696e7465726e616c54797065223a202275696e7438222c0a09090909226e616d65223a20227369676e436f6e6669726d436f756e74222c0a090909092274797065223a202275696e7438220a0909097d2c0a0909097b0a0909090922696e7465726e616c54797065223a2022616464726573735b5d222c0a09090909226e616d65223a20225f616e63686f7273222c0a090909092274797065223a2022616464726573735b5d220a0909097d2c0a0909097b0a0909090922696e7465726e616c54797065223a2022737472696e67222c0a09090909226e616d65223a2022726f75746572222c0a090909092274797065223a2022737472696e67220a0909097d0a09095d2c0a0909226e616d65223a2022636861696e5265676973746572222c0a0909226f757470757473223a205b0a0909097b0a0909090922696e7465726e616c54797065223a2022626f6f6c222c0a09090909226e616d65223a2022222c0a090909092274797065223a2022626f6f6c220a0909097d0a09095d2c0a09092273746174654d75746162696c697479223a20226e6f6e70617961626c65222c0a09092274797065223a202266756e6374696f6e220a097d2c0a097b0a090922696e70757473223a205b0a0909097b0a0909090922636f6d706f6e656e7473223a205b0a09090909097b0a09090909090922696e7465726e616c54797065223a202262797465733332222c0a090909090909226e616d65223a202274784964222c0a0909090909092274797065223a202262797465733332220a09090909097d2c0a09090909097b0a09090909090922696e7465726e616c54797065223a202262797465733332222c0a090909090909226e616d65223a2022747848617368222c0a0909090909092274797065223a202262797465733332220a09090909097d2c0a09090909097b0a09090909090922696e7465726e616c54797065223a2022737472696e67222c0a090909090909226e616d65223a202266726f6d222c0a0909090909092274797065223a2022737472696e67220a09090909097d2c0a09090909097b0a09090909090922696e7465726e616c54797065223a2022737472696e67222c0a090909090909226e616d65223a2022746f222c0a0909090909092274797065223a2022737472696e67220a09090909097d2c0a09090909097b0a09090909090922696e7465726e616c54797065223a2022616464726573732070617961626c65222c0a090909090909226e616d65223a202274616b6572222c0a0909090909092274797065223a202261646472657373220a09090909097d2c0a09090909097b0a09090909090922696e7465726e616c54797065223a202275696e7438222c0a090909090909226e616d65223a20226f726967696e222c0a0909090909092274797065223a202275696e7438220a09090909097d2c0a09090909097b0a09090909090922696e7465726e616c54797065223a202275696e7438222c0a090909090909226e616d65223a2022707572706f7365222c0a0909090909092274797065223a202275696e7438220a09090909097d2c0a09090909097b0a09090909090922696e7465726e616c54797065223a20226279746573222c0a090909090909226e616d65223a202264617461222c0a0909090909092274797065223a20226279746573220a09090909097d0a090909095d2c0a0909090922696e7465726e616c54797065223a20227374727563742043726f73735374727563742e526563657074222c0a09090909226e616d65223a2022727478222c0a090909092274797065223a20227475706c65220a0909097d0a09095d2c0a0909226e616d65223a20226d616b657246696e697368222c0a0909226f757470757473223a205b5d2c0a09092273746174654d75746162696c697479223a202270617961626c65222c0a09092274797065223a202266756e6374696f6e220a097d2c0a097b0a090922616e6f6e796d6f7573223a2066616c73652c0a090922696e70757473223a205b0a0909097b0a0909090922696e6465786564223a2066616c73652c0a0909090922696e7465726e616c54797065223a202262797465733332222c0a09090909226e616d65223a202274784964222c0a090909092274797065223a202262797465733332220a0909097d2c0a0909097b0a0909090922696e6465786564223a2066616c73652c0a0909090922696e7465726e616c54797065223a202261646472657373222c0a09090909226e616d65223a2022746f222c0a090909092274797065223a202261646472657373220a0909097d0a09095d2c0a0909226e616d65223a20224d616b657246696e697368222c0a09092274797065223a20226576656e74220a097d2c0a097b0a090922696e70757473223a205b0a0909097b0a0909090922696e7465726e616c54797065223a202275696e74323536222c0a09090909226e616d65223a20226465737456616c7565222c0a090909092274797065223a202275696e74323536220a0909097d2c0a0909097b0a0909090922696e7465726e616c54797065223a202275696e7438222c0a09090909226e616d65223a2022707572706f7365222c0a090909092274797065223a202275696e7438220a0909097d2c0a0909097b0a0909090922696e7465726e616c54797065223a2022737472696e675b325d222c0a09090909226e616d65223a2022617267222c0a090909092274797065223a2022737472696e675b325d220a0909097d2c0a0909097b0a0909090922696e7465726e616c54797065223a20226279746573222c0a09090909226e616d65223a202264617461222c0a090909092274797065223a20226279746573220a0909097d0a09095d2c0a0909226e616d65223a20226d616b65725374617274222c0a0909226f757470757473223a205b5d2c0a09092273746174654d75746162696c697479223a202270617961626c65222c0a09092274797065223a202266756e6374696f6e220a097d2c0a097b0a090922616e6f6e796d6f7573223a2066616c73652c0a090922696e70757473223a205b0a0909097b0a0909090922696e6465786564223a2066616c73652c0a0909090922696e7465726e616c54797065223a202262797465733332222c0a09090909226e616d65223a202274784964222c0a090909092274797065223a202262797465733332220a0909097d2c0a0909097b0a0909090922696e6465786564223a2066616c73652c0a0909090922696e7465726e616c54797065223a202275696e74323536222c0a09090909226e616d65223a202276616c7565222c0a090909092274797065223a202275696e74323536220a0909097d2c0a0909097b0a0909090922696e6465786564223a2066616c73652c0a0909090922696e7465726e616c54797065223a202275696e74323536222c0a09090909226e616d65223a20226465737456616c7565222c0a090909092274797065223a202275696e74323536220a0909097d2c0a0909097b0a0909090922696e6465786564223a2066616c73652c0a0909090922696e7465726e616c54797065223a2022737472696e67222c0a09090909226e616d65223a202266726f6d222c0a090909092274797065223a2022737472696e67220a0909097d2c0a0909097b0a0909090922696e6465786564223a2066616c73652c0a0909090922696e7465726e616c54797065223a2022737472696e67222c0a09090909226e616d65223a2022746f222c0a090909092274797065223a2022737472696e67220a0909097d2c0a0909097b0a0909090922696e6465786564223a2066616c73652c0a0909090922696e7465726e616c54797065223a202275696e7438222c0a09090909226e616d65223a2022707572706f7365222c0a090909092274797065223a202275696e7438220a0909097d2c0a0909097b0a0909090922696e6465786564223a2066616c73652c0a0909090922696e7465726e616c54797065223a20226279746573222c0a09090909226e616d65223a20227061796c6f6164222c0a090909092274797065223a20226279746573220a0909097d0a09095d2c0a0909226e616d65223a20224d616b65725478222c0a09092274797065223a20226576656e74220a097d2c0a097b0a090922696e70757473223a205b0a0909097b0a0909090922696e7465726e616c54797065223a202275696e7438222c0a09090909226e616d65223a2022707572706f7365222c0a090909092274797065223a202275696e7438220a0909097d2c0a0909097b0a0909090922696e7465726e616c54797065223a2022616464726573735b5d222c0a09090909226e616d65223a20225f616e63686f7273222c0a090909092274797065223a2022616464726573735b5d220a0909097d0a09095d2c0a0909226e616d65223a202272656d6f7665416e63686f7273222c0a0909226f757470757473223a205b5d2c0a09092273746174654d75746162696c697479223a20226e6f6e70617961626c65222c0a09092274797065223a202266756e6374696f6e220a097d2c0a097b0a090922616e6f6e796d6f7573223a2066616c73652c0a090922696e70757473223a205b0a0909097b0a0909090922696e6465786564223a2066616c73652c0a0909090922696e7465726e616c54797065223a202275696e7438222c0a09090909226e616d65223a2022707572706f7365222c0a090909092274797065223a202275696e7438220a0909097d0a09095d2c0a0909226e616d65223a202252656d6f7665416e63686f7273222c0a09092274797065223a20226576656e74220a097d2c0a097b0a090922696e70757473223a205b0a0909097b0a0909090922696e7465726e616c54797065223a202275696e7438222c0a09090909226e616d65223a2022707572706f7365222c0a090909092274797065223a202275696e7438220a0909097d2c0a0909097b0a0909090922696e7465726e616c54797065223a202261646472657373222c0a09090909226e616d65223a20225f616e63686f72222c0a090909092274797065223a202261646472657373220a0909097d2c0a0909097b0a0909090922696e7465726e616c54797065223a2022626f6f6c222c0a09090909226e616d65223a2022737461747573222c0a090909092274797065223a2022626f6f6c220a0909097d0a09095d2c0a0909226e616d65223a2022736574416e63686f72537461747573222c0a0909226f757470757473223a205b5d2c0a09092273746174654d75746162696c697479223a20226e6f6e70617961626c65222c0a09092274797065223a202266756e6374696f6e220a097d2c0a097b0a090922616e6f6e796d6f7573223a2066616c73652c0a090922696e70757473223a205b0a0909097b0a0909090922696e6465786564223a2066616c73652c0a0909090922696e7465726e616c54797065223a202275696e7438222c0a09090909226e616d65223a2022707572706f7365222c0a090909092274797065223a202275696e7438220a0909097d0a09095d2c0a0909226e616d65223a2022536574416e63686f72537461747573222c0a09092274797065223a20226576656e74220a097d2c0a097b0a090922696e70757473223a205b0a0909097b0a0909090922696e7465726e616c54797065223a202275696e7438222c0a09090909226e616d65223a2022707572706f7365222c0a090909092274797065223a202275696e7438220a0909097d2c0a0909097b0a0909090922696e7465726e616c54797065223a202275696e74323536222c0a09090909226e616d65223a20226d617856616c7565222c0a090909092274797065223a202275696e74323536220a0909097d0a09095d2c0a0909226e616d65223a20227365744d617856616c7565222c0a0909226f757470757473223a205b5d2c0a09092273746174654d75746162696c697479223a20226e6f6e70617961626c65222c0a09092274797065223a202266756e6374696f6e220a097d2c0a097b0a090922696e70757473223a205b0a0909097b0a0909090922696e7465726e616c54797065223a202275696e7438222c0a09090909226e616d65223a2022707572706f7365222c0a090909092274797065223a202275696e7438220a0909097d2c0a0909097b0a0909090922696e7465726e616c54797065223a202275696e74323536222c0a09090909226e616d65223a20225f726577617264222c0a090909092274797065223a202275696e74323536220a0909097d0a09095d2c0a0909226e616d65223a2022736574526577617264222c0a0909226f757470757473223a205b5d2c0a09092273746174654d75746162696c697479223a20226e6f6e70617961626c65222c0a09092274797065223a202266756e6374696f6e220a097d2c0a097b0a090922696e70757473223a205b0a0909097b0a0909090922696e7465726e616c54797065223a202275696e7438222c0a09090909226e616d65223a2022707572706f7365222c0a090909092274797065223a202275696e7438220a0909097d2c0a0909097b0a0909090922696e7465726e616c54797065223a202275696e7438222c0a09090909226e616d65223a2022636f756e74222c0a090909092274797065223a202275696e7438220a0909097d0a09095d2c0a0909226e616d65223a20227365745369676e436f6e6669726d436f756e74222c0a0909226f757470757473223a205b5d2c0a09092273746174654d75746162696c697479223a20226e6f6e70617961626c65222c0a09092274797065223a202266756e6374696f6e220a097d2c0a097b0a090922696e70757473223a205b0a0909097b0a0909090922636f6d706f6e656e7473223a205b0a09090909097b0a09090909090922696e7465726e616c54797065223a202262797465733332222c0a090909090909226e616d65223a202274784964222c0a0909090909092274797065223a202262797465733332220a09090909097d2c0a09090909097b0a09090909090922696e7465726e616c54797065223a202262797465733332222c0a090909090909226e616d65223a2022747848617368222c0a0909090909092274797065223a202262797465733332220a09090909097d2c0a09090909097b0a09090909090922696e7465726e616c54797065223a202262797465733332222c0a090909090909226e616d65223a2022626c6f636b48617368222c0a0909090909092274797065223a202262797465733332220a09090909097d2c0a09090909097b0a09090909090922696e7465726e616c54797065223a202275696e74323536222c0a090909090909226e616d65223a202276616c7565222c0a0909090909092274797065223a202275696e74323536220a09090909097d2c0a09090909097b0a09090909090922696e7465726e616c54797065223a202275696e74323536222c0a090909090909226e616d65223a2022636861726765222c0a0909090909092274797065223a202275696e74323536220a09090909097d2c0a09090909097b0a09090909090922696e7465726e616c54797065223a2022616464726573732070617961626c65222c0a090909090909226e616d65223a202266726f6d222c0a0909090909092274797065223a202261646472657373220a09090909097d2c0a09090909097b0a09090909090922696e7465726e616c54797065223a202261646472657373222c0a090909090909226e616d65223a2022746f222c0a0909090909092274797065223a202261646472657373220a09090909097d2c0a09090909097b0a09090909090922696e7465726e616c54797065223a202275696e7438222c0a090909090909226e616d65223a20226f726967696e222c0a0909090909092274797065223a202275696e7438220a09090909097d2c0a09090909097b0a09090909090922696e7465726e616c54797065223a202275696e7438222c0a090909090909226e616d65223a2022707572706f7365222c0a0909090909092274797065223a202275696e7438220a09090909097d2c0a09090909097b0a09090909090922696e7465726e616c54797065223a20226279746573222c0a090909090909226e616d65223a20227061796c6f6164222c0a0909090909092274797065223a20226279746573220a09090909097d2c0a09090909097b0a09090909090922696e7465726e616c54797065223a202275696e743235365b5d222c0a090909090909226e616d65223a202276222c0a0909090909092274797065223a202275696e743235365b5d220a09090909097d2c0a09090909097b0a09090909090922696e7465726e616c54797065223a2022627974657333325b5d222c0a090909090909226e616d65223a202272222c0a0909090909092274797065223a2022627974657333325b5d220a09090909097d2c0a09090909097b0a09090909090922696e7465726e616c54797065223a2022627974657333325b5d222c0a090909090909226e616d65223a202273222c0a0909090909092274797065223a2022627974657333325b5d220a09090909097d0a090909095d2c0a0909090922696e7465726e616c54797065223a20227374727563742043726f73735374727563742e4f72646572222c0a09090909226e616d65223a2022637478222c0a090909092274797065223a20227475706c65220a0909097d2c0a0909097b0a0909090922696e7465726e616c54797065223a2022737472696e67222c0a09090909226e616d65223a2022746f222c0a090909092274797065223a2022737472696e67220a0909097d2c0a0909097b0a0909090922696e7465726e616c54797065223a20226279746573222c0a09090909226e616d65223a202264617461222c0a090909092274797065223a20226279746573220a0909097d0a09095d2c0a0909226e616d65223a202274616b6572222c0a0909226f757470757473223a205b5d2c0a09092273746174654d75746162696c697479223a202270617961626c65222c0a09092274797065223a202266756e6374696f6e220a097d2c0a097b0a090922616e6f6e796d6f7573223a2066616c73652c0a090922696e70757473223a205b0a0909097b0a0909090922696e6465786564223a2066616c73652c0a0909090922696e7465726e616c54797065223a202262797465733332222c0a09090909226e616d65223a202274784964222c0a090909092274797065223a202262797465733332220a0909097d2c0a0909097b0a0909090922696e6465786564223a2066616c73652c0a0909090922696e7465726e616c54797065223a202261646472657373222c0a09090909226e616d65223a202266726f6d222c0a090909092274797065223a202261646472657373220a0909097d2c0a0909097b0a0909090922696e6465786564223a2066616c73652c0a0909090922696e7465726e616c54797065223a2022737472696e67222c0a09090909226e616d65223a2022746f222c0a090909092274797065223a2022737472696e67220a0909097d2c0a0909097b0a0909090922696e6465786564223a2066616c73652c0a0909090922696e7465726e616c54797065223a202275696e7438222c0a09090909226e616d65223a2022636861696e222c0a090909092274797065223a202275696e7438220a0909097d2c0a0909097b0a0909090922696e6465786564223a2066616c73652c0a0909090922696e7465726e616c54797065223a20226279746573222c0a09090909226e616d65223a20227061796c6f6164222c0a090909092274797065223a20226279746573220a0909097d0a09095d2c0a0909226e616d65223a202254616b65725478222c0a09092274797065223a20226576656e74220a097d2c0a097b0a090922696e70757473223a205b0a0909097b0a0909090922696e7465726e616c54797065223a202275696e7438222c0a09090909226e616d65223a2022707572706f7365222c0a090909092274797065223a202275696e7438220a0909097d2c0a0909097b0a0909090922696e7465726e616c54797065223a2022737472696e67222c0a09090909226e616d65223a20225f726f75746572222c0a090909092274797065223a2022737472696e67220a0909097d0a09095d2c0a0909226e616d65223a2022757064617465526f75746572222c0a0909226f757470757473223a205b5d2c0a09092273746174654d75746162696c697479223a20226e6f6e70617961626c65222c0a09092274797065223a202266756e6374696f6e220a097d2c0a097b0a090922616e6f6e796d6f7573223a2066616c73652c0a090922696e70757473223a205b0a0909097b0a0909090922696e6465786564223a2066616c73652c0a0909090922696e7465726e616c54797065223a202275696e7438222c0a09090909226e616d65223a2022707572706f7365222c0a090909092274797065223a202275696e7438220a0909097d0a09095d2c0a0909226e616d65223a2022557064617465526f75746572222c0a09092274797065223a20226576656e74220a097d2c0a097b0a090922696e70757473223a205b0a0909097b0a0909090922696e7465726e616c54797065223a202275696e743634222c0a09090909226e616d65223a20226e222c0a090909092274797065223a202275696e743634220a0909097d0a09095d2c0a0909226e616d65223a2022626974436f756e74222c0a0909226f757470757473223a205b0a0909097b0a0909090922696e7465726e616c54797065223a202275696e743634222c0a09090909226e616d65223a2022222c0a090909092274797065223a202275696e743634220a0909097d0a09095d2c0a09092273746174654d75746162696c697479223a202270757265222c0a09092274797065223a202266756e6374696f6e220a097d2c0a097b0a090922696e70757473223a205b5d2c0a0909226e616d65223a2022636861696e4964222c0a0909226f757470757473223a205b0a0909097b0a0909090922696e7465726e616c54797065223a202275696e74323536222c0a09090909226e616d65223a20226964222c0a090909092274797065223a202275696e74323536220a0909097d0a09095d2c0a09092273746174654d75746162696c697479223a202270757265222c0a09092274797065223a202266756e6374696f6e220a097d2c0a097b0a090922696e70757473223a205b0a0909097b0a0909090922696e7465726e616c54797065223a202275696e7438222c0a09090909226e616d65223a2022222c0a090909092274797065223a202275696e7438220a0909097d0a09095d2c0a0909226e616d65223a202263726f7373436861696e73222c0a0909226f757470757473223a205b0a0909097b0a0909090922696e7465726e616c54797065223a202275696e7438222c0a09090909226e616d65223a2022707572706f7365222c0a090909092274797065223a202275696e7438220a0909097d2c0a0909097b0a0909090922696e7465726e616c54797065223a202275696e7438222c0a09090909226e616d65223a20227369676e436f6e6669726d436f756e74222c0a090909092274797065223a202275696e7438220a0909097d2c0a0909097b0a0909090922696e7465726e616c54797065223a202275696e74323536222c0a09090909226e616d65223a20226d617856616c7565222c0a090909092274797065223a202275696e74323536220a0909097d2c0a0909097b0a0909090922696e7465726e616c54797065223a202275696e743634222c0a09090909226e616d65223a2022616e63686f7273506f736974696f6e426974222c0a090909092274797065223a202275696e743634220a0909097d2c0a0909097b0a0909090922696e7465726e616c54797065223a202275696e743634222c0a09090909226e616d65223a202264656c73506f736974696f6e426974222c0a090909092274797065223a202275696e743634220a0909097d2c0a0909097b0a0909090922696e7465726e616c54797065223a202275696e7438222c0a09090909226e616d65223a202264656c4964222c0a090909092274797065223a202275696e7438220a0909097d2c0a0909097b0a0909090922696e7465726e616c54797065223a202275696e74323536222c0a09090909226e616d65223a2022726577617264222c0a090909092274797065223a202275696e74323536220a0909097d2c0a0909097b0a0909090922696e7465726e616c54797065223a202275696e74323536222c0a09090909226e616d65223a2022746f74616c526577617264222c0a090909092274797065223a202275696e74323536220a0909097d2c0a0909097b0a0909090922696e7465726e616c54797065223a2022737472696e67222c0a09090909226e616d65223a2022726f75746572222c0a090909092274797065223a2022737472696e67220a0909097d0a09095d2c0a09092273746174654d75746162696c697479223a202276696577222c0a09092274797065223a202266756e6374696f6e220a097d2c0a097b0a090922696e70757473223a205b0a0909097b0a0909090922696e7465726e616c54797065223a202275696e7438222c0a09090909226e616d65223a2022707572706f7365222c0a090909092274797065223a202275696e7438220a0909097d0a09095d2c0a0909226e616d65223a2022676574416e63686f7273222c0a0909226f757470757473223a205b0a0909097b0a0909090922696e7465726e616c54797065223a2022616464726573735b5d222c0a09090909226e616d65223a20225f616e63686f7273222c0a090909092274797065223a2022616464726573735b5d220a0909097d2c0a0909097b0a0909090922696e7465726e616c54797065223a202275696e7438222c0a09090909226e616d65223a2022222c0a090909092274797065223a202275696e7438220a0909097d0a09095d2c0a09092273746174654d75746162696c697479223a202276696577222c0a09092274797065223a202266756e6374696f6e220a097d2c0a097b0a090922696e70757473223a205b0a0909097b0a0909090922696e7465726e616c54797065223a202275696e7438222c0a09090909226e616d65223a2022707572706f7365222c0a090909092274797065223a202275696e7438220a0909097d2c0a0909097b0a0909090922696e7465726e616c54797065223a202261646472657373222c0a09090909226e616d65223a20225f616e63686f72222c0a090909092274797065223a202261646472657373220a0909097d0a09095d2c0a0909226e616d65223a2022676574416e63686f72576f726b436f756e74222c0a0909226f757470757473223a205b0a0909097b0a0909090922696e7465726e616c54797065223a202275696e74323536222c0a09090909226e616d65223a2022222c0a090909092274797065223a202275696e74323536220a0909097d2c0a0909097b0a0909090922696e7465726e616c54797065223a202275696e74323536222c0a09090909226e616d65223a2022222c0a090909092274797065223a202275696e74323536220a0909097d0a09095d2c0a09092273746174654d75746162696c697479223a202276696577222c0a09092274797065223a202266756e6374696f6e220a097d2c0a097b0a090922696e70757473223a205b0a0909097b0a0909090922696e7465726e616c54797065223a202275696e7438222c0a09090909226e616d65223a2022707572706f7365222c0a090909092274797065223a202275696e7438220a0909097d0a09095d2c0a0909226e616d65223a2022676574436861696e526577617264222c0a0909226f757470757473223a205b0a0909097b0a0909090922696e7465726e616c54797065223a202275696e74323536222c0a09090909226e616d65223a2022222c0a090909092274797065223a202275696e74323536220a0909097d0a09095d2c0a09092273746174654d75746162696c697479223a202276696577222c0a09092274797065223a202266756e6374696f6e220a097d2c0a097b0a090922696e70757473223a205b0a0909097b0a0909090922696e7465726e616c54797065223a202275696e7438222c0a09090909226e616d65223a2022707572706f7365222c0a090909092274797065223a202275696e7438220a0909097d2c0a0909097b0a0909090922696e7465726e616c54797065223a202261646472657373222c0a09090909226e616d65223a20225f616e63686f72222c0a090909092274797065223a202261646472657373220a0909097d0a09095d2c0a0909226e616d65223a202267657444656c416e63686f725369676e436f756e74222c0a0909226f757470757473223a205b0a0909097b0a0909090922696e7465726e616c54797065223a202275696e74323536222c0a09090909226e616d65223a2022222c0a090909092274797065223a202275696e74323536220a0909097d0a09095d2c0a09092273746174654d75746162696c697479223a202276696577222c0a09092274797065223a202266756e6374696f6e220a097d2c0a097b0a090922696e70757473223a205b0a0909097b0a0909090922696e7465726e616c54797065223a202262797465733332222c0a09090909226e616d65223a202274784964222c0a090909092274797065223a202262797465733332220a0909097d2c0a0909097b0a0909090922696e7465726e616c54797065223a202275696e7438222c0a09090909226e616d65223a2022707572706f7365222c0a090909092274797065223a202275696e7438220a0909097d0a09095d2c0a0909226e616d65223a20226765744d616b65725478222c0a0909226f757470757473223a205b0a0909097b0a0909090922696e7465726e616c54797065223a202275696e74323536222c0a09090909226e616d65223a2022222c0a090909092274797065223a202275696e74323536220a0909097d0a09095d2c0a09092273746174654d75746162696c697479223a202276696577222c0a09092274797065223a202266756e6374696f6e220a097d2c0a097b0a090922696e70757473223a205b0a0909097b0a0909090922696e7465726e616c54797065223a202275696e7438222c0a09090909226e616d65223a2022707572706f7365222c0a090909092274797065223a202275696e7438220a0909097d0a09095d2c0a0909226e616d65223a20226765744d617856616c7565222c0a0909226f757470757473223a205b0a0909097b0a0909090922696e7465726e616c54797065223a202275696e74323536222c0a09090909226e616d65223a2022222c0a090909092274797065223a202275696e74323536220a0909097d0a09095d2c0a09092273746174654d75746162696c697479223a202276696577222c0a09092274797065223a202266756e6374696f6e220a097d2c0a097b0a090922696e70757473223a205b0a0909097b0a0909090922696e7465726e616c54797065223a202275696e7438222c0a09090909226e616d65223a2022707572706f7365222c0a090909092274797065223a202275696e7438220a0909097d0a09095d2c0a0909226e616d65223a2022676574526f75746572222c0a0909226f757470757473223a205b0a0909097b0a0909090922696e7465726e616c54797065223a2022737472696e67222c0a09090909226e616d65223a2022222c0a090909092274797065223a2022737472696e67220a0909097d0a09095d2c0a09092273746174654d75746162696c697479223a202276696577222c0a09092274797065223a202266756e6374696f6e220a097d2c0a097b0a090922696e70757473223a205b0a0909097b0a0909090922696e7465726e616c54797065223a202262797465733332222c0a09090909226e616d65223a202274784964222c0a090909092274797065223a202262797465733332220a0909097d2c0a0909097b0a0909090922696e7465726e616c54797065223a202261646472657373222c0a09090909226e616d65223a20225f66726f6d222c0a090909092274797065223a202261646472657373220a0909097d2c0a0909097b0a0909090922696e7465726e616c54797065223a202275696e7438222c0a09090909226e616d65223a2022707572706f7365222c0a090909092274797065223a202275696e7438220a0909097d0a09095d2c0a0909226e616d65223a202267657454616b65725478222c0a0909226f757470757473223a205b0a0909097b0a0909090922696e7465726e616c54797065223a202275696e74323536222c0a09090909226e616d65223a2022222c0a090909092274797065223a202275696e74323536220a0909097d0a09095d2c0a09092273746174654d75746162696c697479223a202276696577222c0a09092274797065223a202266756e6374696f6e220a097d2c0a097b0a090922696e70757473223a205b0a0909097b0a0909090922696e7465726e616c54797065223a202275696e7438222c0a09090909226e616d65223a2022707572706f7365222c0a090909092274797065223a202275696e7438220a0909097d0a09095d2c0a0909226e616d65223a2022676574546f74616c526577617264222c0a0909226f757470757473223a205b0a0909097b0a0909090922696e7465726e616c54797065223a202275696e74323536222c0a09090909226e616d65223a2022222c0a090909092274797065223a202275696e74323536220a0909097d0a09095d2c0a09092273746174654d75746162696c697479223a202276696577222c0a09092274797065223a202266756e6374696f6e220a097d2c0a097b0a090922696e70757473223a205b5d2c0a0909226e616d65223a20226c697374222c0a0909226f757470757473223a205b0a0909097b0a0909090922696e7465726e616c54797065223a202275696e74323536222c0a09090909226e616d65223a20226c6c222c0a090909092274797065223a202275696e74323536220a0909097d0a09095d2c0a09092273746174654d75746162696c697479223a202270757265222c0a09092274797065223a202266756e6374696f6e220a097d2c0a097b0a090922696e70757473223a205b5d2c0a0909226e616d65223a20226f776e6572222c0a0909226f757470757473223a205b0a0909097b0a0909090922696e7465726e616c54797065223a202261646472657373222c0a09090909226e616d65223a2022222c0a090909092274797065223a202261646472657373220a0909097d0a09095d2c0a09092273746174654d75746162696c697479223a202276696577222c0a09092274797065223a202266756e6374696f6e220a097d0a5d"

func init()  {
	data, err := hexutil.Decode(CrossAbi)
	if err != nil {
		log.Error("init","err",err)
	}
	abiParsed, err = abi.JSON(bytes.NewReader(data))
	if err != nil {
		log.Error("abi.JSON","err",err)
	}
}

type Viewer struct {
	Client       	*rpc.Client
	SimpleClient 	*ethclient.Client
	Address      	string
	currentHeight   uint64
	ctmsChan     	chan *core.CrossTransaction
	PrivateKey      *ecdsa.PrivateKey

	ctx    context.Context
	cancel context.CancelFunc
}

func New(repo *repo.Repo,ch chan *core.CrossTransaction) (*Viewer,error) {
	ctx, cancel := context.WithCancel(context.Background())
	//log.Info("New","addr",repo.Config.RpcUrl)
	client, err := rpc.DialContext(ctx,fmt.Sprintf("http://%s:%s", repo.Config.RpcIp, repo.Config.RpcPort))
	if err != nil {
		return nil, err
	}
	return &Viewer{
		Client: client,
		SimpleClient: ethclient.NewClient(client),
		Address: repo.Config.Contract,
		currentHeight: 3937000,
		ctmsChan: ch,
		PrivateKey: repo.Key.PrivKey.(*ecdsa.PrivateKey),
		ctx: ctx,
		cancel: cancel,
	},nil
}

func (this *Viewer)Start() error {
	go this.loop()
	return nil
}

func (this *Viewer)Stop() error {
	this.cancel()
	return nil
}

func (this *Viewer)loop()  {
	var eventTicker = time.NewTicker(time.Second*5)
	defer eventTicker.Stop()
	for {
		select {
		case <-this.ctx.Done():
			return
		case <-eventTicker.C:
			this.GetEvents()
		}
	}
}

func (this *Viewer)GetEvents() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	var result hexutil.Big
	if err := this.Client.CallContext(ctx,&result, "eth_blockNumber");err != nil {
		log.Info("CallContext","err",err)
	}
	var toBlock uint64
	if this.currentHeight < result.ToInt().Uint64() - 100 {
		toBlock = this.currentHeight + 100
	} else {
		toBlock = result.ToInt().Uint64() - 12
	}
	records := simplechain.FilterQuery{
		FromBlock: big.NewInt(int64(this.currentHeight)),
		ToBlock: big.NewInt(int64(toBlock)),
		Addresses: []common.Address{common.HexToAddress(this.Address)},
	}
	logs, err := this.SimpleClient.FilterLogs(ctx, records)
	if err != nil {
		log.Info("GetEvents","err",err)
		//TODO err
	}
	if len(logs) > 0 {
		this.EventLog(logs)
	}
	this.currentHeight = toBlock
	log.Info("GetEvents","currentHeight",this.currentHeight)

}


func (this *Viewer) EventLog(logs []types.Log) {
	makerTx := abiParsed.Events["MakerTx"].ID().Hex()
	takerTx := abiParsed.Events["TakerTx"].ID().Hex()
	makerFinish := abiParsed.Events["MakerFinish"].ID().Hex()
	for _, event := range logs {
		switch event.Topics[0].Hex() {
		case makerTx:
			var args CrossMakerTx
			err := abiParsed.Unpack(&args, "MakerTx", event.Data)
			if err != nil {
				log.Info("EventLog","Unpack err",err)
			}

			ctm :=  core.NewCrossTransaction(args.Value,args.DestValue,args.From,args.To,1,args.Purpose, args.TxId,event.TxHash,event.BlockHash,args.Payload)
			signHash := func(hash []byte) ([]byte, error) {
				return  crypto.Sign(hash,this.PrivateKey.K)
			}
			ctms,err :=  core.SignCtx(ctm,core.MakeCtxSigner(big.NewInt(11)),signHash)
			if err != nil {
				log.Info("SignCtx","err",err)
			}
			from,err := core.CtxSender(core.MakeCtxSigner(big.NewInt(11)),ctms)
			if err != nil {
				log.Info("CtxSender","err",err)
			}
			publicKey :=  crypto.PubkeyToAddress(this.PrivateKey.K.PublicKey)
			if err != nil {
				log.Info("PublicKey","err",err)
			}
			log.Info("receive ctx msg","msg",args,"ctms",ctms,"sender",from.String(),"publicKey",publicKey.String())
			this.ctmsChan <- ctms
			//TODO ctx and Handle
			//ctxMsg := core.NewCrossTransaction(args.Value,args.DestValue,args.From,args.To,)
		case takerTx:
			var args CrossTakerTx
			err := abiParsed.Unpack(&args, "TakerTx", event.Data)
			if err != nil {
				log.Info("EventLog","Unpack err",err)
			}
			log.Info("receive rtx msg","msg",args.RemoteChainId.String())
			//TODO rtx and handle
			//rtxMsg := core.NewReceptTransaction()
		case makerFinish:
			log.Info("receive finish msg")
			//TODO finish handle
		}
	}
}

type CrossMakerTx struct {
	TxId          [32]byte
	Value         *big.Int
	DestValue     *big.Int
	From          string
	To            string
	Purpose       uint8
	Payload       []byte
	//Raw           types.Log
}

type CrossTakerTx struct {
	TxId          [32]byte
	To            common.Address
	RemoteChainId *big.Int
	From          common.Address
	Value         *big.Int
	DestValue     *big.Int
	//Raw           types.Log
}

type CrossMakerFinish struct {
	TxId [32]byte
	To   common.Address
	//Raw  types.Log
}

// 实现区块链交互功能。
// 签署消息、发送到地址、设置交易手续费、发送RPC等四个函数。
package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/bitspill/flojson" // FLO JSON-RPC API
)

var (
	id     int64
	user   string
	pass   string
	server string
)

func init() {
	id = 0                                 // id is static at 0, for "proper" json-rpc increment with each call
	user = config.FloConfiguration.RpcUser // 读取 config.go 中的设置
	pass = config.FloConfiguration.RpcPass
	server = config.FloConfiguration.RpcAddress
}

func signMessage(address, message string) (string, error) { // 签署信息
	//fmt.Println(id) //0
	//fmt.Println(address) //FS267EyRkAEyLNuSaTdbtkgcjjsXEBYhse
	//fmt.Println(message) //QmfJxwEBCbfe5SRQpP1T1jaJrCLMuBSwP8FgpeV54rsPLx-FS267EyRkAEyLNuSaTdbtkgcjjsXEBYhse-1604911864
	cmd, err := flojson.NewSignMessageCmd(id, address, message)
	//fmt.Println(cmd) //&{0 FS267EyRkAEyLNuSaTdbtkgcjjsXEBYhse QmfJxwEBCbfe5SRQpP1T1jaJrCLMuBSwP8FgpeV54rsPLx-FS267EyRkAEyLNuSaTdbtkgcjjsXEBYhse-1604911932}
	//fmt.Println(err) //<nil>
	if err != nil {
		return "", err
	}

	reply, err := sendRPC(cmd)
	//fmt.Println(reply) //Error getting json reply: Error sending json message: Post http://<username>:<password>@localhost:17317: dial tcp 127.0.0.1:17317: connect: connection refused
	//fmt.Println(err) //{<nil> <nil> <nil>}
	if err != nil {
		return "", err
	}

	if signature, ok := reply.Result.(string); ok {
		return signature, nil
	}

	return "", errors.New("unexpected rpc error")
}

func sendToAddress(address string, amount float64, floData string) (string, error) { // 发送信息
	satoshi := int64(1e8 * amount)
	cmd, err := flojson.NewSendToAddressCmd(id, address, satoshi, "", "", floData)
	if err != nil {
		return "", err
	}

	reply, err := sendRPC(cmd)
	if err != nil {
		return "", err
	}
	if reply.Error != nil {
		return "", reply.Error
	}
	return reply.Result.(string), nil
}

func setTxFee(floPerKb float64) error { // 设置交易费用
	var satoshi = int64(floPerKb * 1e8)
	cmd, err := flojson.NewSetTxFeeCmd(id, satoshi)
	if err != nil {
		return err
	}

	reply, err := sendRPC(cmd)
	if err != nil {
		return err
	}
	if reply.Error != nil {
		return reply.Error
	}
	return nil
}

// 因为测试而改写的程序
// 因为测试而改写的程序
// 因为测试而改写的程序

// Cmd is an interface for all Bitcoin JSON API commands to marshal
// and unmarshal as a JSON object.
type Cmd interface {
	json.Marshaler
	json.Unmarshaler
	Id() interface{}
	Method() string
}

// RawCmd is a type for unmarshaling raw commands into before the
// custom command type is set.  Other packages may register their
// own RawCmd to Cmd converters by calling RegisterCustomCmd.
type RawCmd struct {
	Jsonrpc string            `json:"jsonrpc"`
	Id      interface{}       `json:"id"`
	Method  string            `json:"method"`
	Params  []json.RawMessage `json:"params"`
}

// RawCmdParser is a function to create a custom Cmd from a RawCmd.
type RawCmdParser func(*RawCmd) (Cmd, error)

// ReplyParser is a function a custom Cmd can use to unmarshal the results of a
// reply into a concrete struct.
type ReplyParser func(json.RawMessage) (interface{}, error)

type cmd struct {
	parser      RawCmdParser
	replyParser ReplyParser
	helpString  string
}

var customCmds = make(map[string]cmd)

// RpcSend sends the passed command to the provided server using the provided
// authentication details, waits for a reply, and returns a Go struct with the
// result.
func RpcSend(user string, password string, server string, cmd flojson.Cmd) (flojson.Reply, error) {
	msg, err := cmd.MarshalJSON()
	//fmt.Println(msg) //[123 34 106 115 111 110 114 112 99 34 58 34 49 46 48 34 44 34 105 100 34 58 48 44 34 109 101 116 104 111 100 34 58 34 115 105 103 110 109 101 115 115 97 103 101 34 44 34 112 97 114 97 109 115 34 58 91 34 70 83 50 54 55 69 121 82 107 65 69 121 76 78 117 83 97 84 100 98 116 107 103 99 106 106 115 88 69 66 89 104 115 101 34 44 34 81 109 102 74 120 119 69 66 67 98 102 101 53 83 82 81 112 80 49 84 49 106 97 74 114 67 76 77 117 66 83 119 80 56 70 103 112 101 86 53 52 114 115 80 76 120 45 70 83 50 54 55 69 121 82 107 65 69 121 76 78 117 83 97 84 100 98 116 107 103 99 106 106 115 88 69 66 89 104 115 101 45 49 54 48 52 57 50 55 50 55 54 34 93 125]
	if err != nil {
		return flojson.Reply{}, err
	}

	return RpcCommand(user, password, server, msg)
}

// RpcCommand takes a message generated from one of the routines above
// along with the login/server info, sends it, and gets a reply, returning
// a go struct with the result.
func RpcCommand(user string, password string, server string, message []byte) (flojson.Reply, error) {
	return rpcCommand(user, password, server, message, false, nil, false)
}

func rpcCommand(user string, password string, server string, message []byte,
	https bool, certificates []byte, skipverify bool) (flojson.Reply, error) {
	var result flojson.Reply
	//fmt.Println(result) //{<nil> <nil> <nil>}
	method, err := flojson.JSONGetMethod(message)
	//fmt.Println(method) //signmessage
	//fmt.Println(err) //<nil>
	if err != nil {
		return result, err
	}
	//fmt.Println(user) //user
	//fmt.Println(password) //pass
	//fmt.Println(server) //localhost:17317
	//fmt.Println(message) //[123 34 106 115 111 110 114 112 99 34 58 34 49 46 48 34 44 34 105 100 34 58 48 44 34 109 101 116 104 111 100 34 58 34 115 105 103 110 109 101 115 115 97 103 101 34 44 34 112 97 114 97 109 115 34 58 91 34 70 83 50 54 55 69 121 82 107 65 69 121 76 78 117 83 97 84 100 98 116 107 103 99 106 106 115 88 69 66 89 104 115 101 34 44 34 81 109 102 74 120 119 69 66 67 98 102 101 53 83 82 81 112 80 49 84 49 106 97 74 114 67 76 77 117 66 83 119 80 56 70 103 112 101 86 53 52 114 115 80 76 120 45 70 83 50 54 55 69 121 82 107 65 69 121 76 78 117 83 97 84 100 98 116 107 103 99 106 106 115 88 69 66 89 104 115 101 45 49 54 48 52 57 55 56 56 56 54 34 93 125]
	//fmt.Println(https) //false
	//fmt.Println(certificates) //[]
	//fmt.Println(skipverify) //false
	body, err := rpcRawCommand(user, password, server, message, https, certificates, skipverify)
	//body, err := flojson.RpcRawCommand(user, password, server, message)
	//fmt.Println(body) //[67 108 105 101 110 116 32 115 101 110 116 32 97 110 32 72 84 84 80 32 114 101 113 117 101 115 116 32 116 111 32 97 110 32 72 84 84 80 83 32 115 101 114 118 101 114 46 10]
	//fmt.Println(err) //<nil>
	if err != nil {
		err := fmt.Errorf("Error getting json reply: %v", err)
		return result, err
	}
	//fmt.Println(method) //signmessage
	//fmt.Println(body) //[67 108 105 101 110 116 32 115 101 110 116 32 97 110 32 72 84 84 80 32 114 101 113 117 101 115 116 32 116 111 32 97 110 32 72 84 84 80 83 32 115 101 114 118 101 114 46 10]
	result, err = ReadResultCmd(method, body)
	//fmt.Println(result) //{<nil> <nil> <nil>}
	//fmt.Println(result) //{<nil> <nil> <nil>}
	//fmt.Println(result) //{<nil> <nil> <nil>}
	//fmt.Println(err) //Error unmarshalling json reply: invalid character 'C' looking for beginning of value
	//fmt.Println(err)
	//fmt.Println(err)
	if err != nil {
		err := fmt.Errorf("Error reading json message: %v", err)
		return result, err
	}
	return result, err
}

// rpcRawCommand is a helper function for the above two functions.
func rpcRawCommand(user string, password string, server string,
	message []byte, https bool, certificates []byte, skipverify bool) ([]byte, error) {
	var result []byte
	var msg interface{}
	err := json.Unmarshal(message, &msg)
	//fmt.Println(msg) //map[jsonrpc:1.0 id:0 method:signmessage params:[FS267EyRkAEyLNuSaTdbtkgcjjsXEBYhse QmfJxwEBCbfe5SRQpP1T1jaJrCLMuBSwP8FgpeV54rsPLx-FS267EyRkAEyLNuSaTdbtkgcjjsXEBYhse-1604993249]]
	//fmt.Println(err) //<nil>
	if err != nil {
		err := fmt.Errorf("Error, message does not appear to be valid json: %v", err)
		return result, err
	}
	//fmt.Println(user) //user
	//fmt.Println(password) //pass
	//fmt.Println(server) //localhost:17317
	//fmt.Println(message) //[123 34 106 115 111 110 114 112 99 34 58 34 49 46 48 34 44 34 105 100 34 58 48 44 34 109 101 116 104 111 100 34 58 34 115 105 103 110 109 101 115 115 97 103 101 34 44 34 112 97 114 97 109 115 34 58 91 34 70 83 50 54 55 69 121 82 107 65 69 121 76 78 117 83 97 84 100 98 116 107 103 99 106 106 115 88 69 66 89 104 115 101 34 44 34 81 109 102 74 120 119 69 66 67 98 102 101 53 83 82 81 112 80 49 84 49 106 97 74 114 67 76 77 117 66 83 119 80 56 70 103 112 101 86 53 52 114 115 80 76 120 45 70 83 50 54 55 69 121 82 107 65 69 121 76 78 117 83 97 84 100 98 116 107 103 99 106 106 115 88 69 66 89 104 115 101 45 49 54 48 52 57 55 56 56 56 54 34 93 125]
	//fmt.Println(https) //false
	//fmt.Println(certificates) //[]
	//fmt.Println(skipverify) //false
	resp, err := jsonRpcSend(user, password, server, message, true, certificates, true)
	//fmt.Println(resp) //&{400 Bad Request 400 HTTP/1.0 1 0 map[] 0xc420220240 -1 [] true false map[] 0xc420024500 <nil>}
	//fmt.Println(resp.Body) // &{0xc4200ac500 {0 0} false <nil> 0x6574a0 0x657430}
	//fmt.Println(err) //<nil>
	if err != nil {
		err := fmt.Errorf("Error sending json message: " + err.Error())
		return result, err
	}
	result, err = flojson.GetRaw(resp.Body)
	fmt.Println(result) //[67 108 105 101 110 116 32 115 101 110 116 32 97 110 32 72 84 84 80 32 114 101 113 117 101 115 116 32 116 111 32 97 110 32 72 84 84 80 83 32 115 101 114 118 101 114 46 10]
	// 修改后变成[123 34 114 101 115 117 108 116 34 58 110 117 108 108 44 34 101 114 114 111 114 34 58 123 34 99 111 100 101 34 58 45 49 44 34 109 101 115 115 97 103 101 34 58 34 84 104 105 115 32 105 109 112 108 101 109 101 110 116 97 116 105 111 110 32 100 111 101 115 32 110 111 116 32 105 109 112 108 101 109 101 110 116 32 119 97 108 108 101 116 32 99 111 109 109 97 110 100 115 34 125 44 34 105 100 34 58 48 125 10]
	fmt.Println(err) //<nil>
	if err != nil {
		err := fmt.Errorf("Error getting json reply: %v", err)
		return result, err
	}
	return result, err
}

// jsonRpcSend connects to the daemon with the specified username, password,
// and ip/port and then send the supplied message.  This uses net/http rather
// than net/rpc/jsonrpc since that one doesn't support http connections and is
// therefore useless.
func jsonRpcSend(user string, password string, server string, message []byte,
	https bool, certificates []byte, skipverify bool) (*http.Response, error) {
	client := &http.Client{}
	protocol := "http"
	if https {
		pool := x509.NewCertPool()
		pool.AppendCertsFromPEM(certificates)

		config := &tls.Config{
			InsecureSkipVerify: skipverify,
			RootCAs:            pool,
		}
		transport := &http.Transport{TLSClientConfig: config}
		client.Transport = transport
		protocol = "https"
	}
	user = url.PathEscape(user)
	//fmt.Println(user) //user
	password = url.PathEscape(password)
	//fmt.Println(password) //pass
	credentials := user + ":" + password
	//fmt.Println(credentials) //user:pass
	//fmt.Println(bytes.NewReader(message)) //&{[123 34 106 115 111 110 114 112 99 34 58 34 49 46 48 34 44 34 105 100 34 58 48 44 34 109 101 116 104 111 100 34 58 34 115 105 103 110 109 101 115 115 97 103 101 34 44 34 112 97 114 97 109 115 34 58 91 34 70 83 50 54 55 69 121 82 107 65 69 121 76 78 117 83 97 84 100 98 116 107 103 99 106 106 115 88 69 66 89 104 115 101 34 44 34 81 109 102 74 120 119 69 66 67 98 102 101 53 83 82 81 112 80 49 84 49 106 97 74 114 67 76 77 117 66 83 119 80 56 70 103 112 101 86 53 52 114 115 80 76 120 45 70 83 50 54 55 69 121 82 107 65 69 121 76 78 117 83 97 84 100 98 116 107 103 99 106 106 115 88 69 66 89 104 115 101 45 49 54 48 52 57 57 51 56 49 51 34 93 125] 0 -1}
	resp, err := client.Post(protocol+"://"+credentials+"@"+server,
		"application/json", bytes.NewReader(message))
	//fmt.Println(resp) //&{400 Bad Request 400 HTTP/1.0 1 0 map[] 0xc4202321c0 -1 [] true false map[] 0xc420024700 <nil>}
	//修改后&{200 OK 200 HTTP/1.1 1 1 map[Content-Type:[application/json]] 0xc4201fc180 -1 [] true false map[] 0xc420262300 0xc4200ec000}
	//fmt.Println(err) //<nil>
	if err != nil {
		// We do not want to log the username/password in the errors.
		replaceStr := "<username>:<password>"
		str := strings.Replace(err.Error(), credentials, replaceStr, -1)
		err = fmt.Errorf("%v", str)
	}
	return resp, err
}

//
//// decodeState represents the state while decoding a JSON value.
//type decodeState struct {
//	data         []byte
//	off          int // read offset in data
//	scan         scanner
//	nextscan     scanner  // for calls to nextValue
//	errorContext struct { // provides context for type errors
//		Struct string
//		Field  string
//	}
//	savedError            error
//	useNumber             bool
//	disallowUnknownFields bool
//}
//
//type scanner struct {
//	// The step is a func to be called to execute the next transition.
//	// Also tried using an integer constant and a single func
//	// with a switch, but using the func directly was 10% faster
//	// on a 64-bit Mac Mini, and it's nicer to read.
//	step func(*scanner, byte) int
//
//	// Reached end of top-level value.
//	endTop bool
//
//	// Stack of what we're in the middle of - array values, object keys, object values.
//	parseState []int
//
//	// Error that happened, if any.
//	err error
//
//	// 1-byte redo (see undo method)
//	redo      bool
//	redoCode  int
//	redoState func(*scanner, byte) int
//
//	// total bytes consumed, updated by decoder.Decode
//	bytes int64
//}
//
//// reset prepares the scanner for use.
//// It must be called before calling s.step.
//func (s *scanner) reset() {
//	s.step = stateBeginValue
//	s.parseState = s.parseState[0:0]
//	s.err = nil
//	s.redo = false
//	s.endTop = false
//}
//
//
//// eof tells the scanner that the end of input has been reached.
//// It returns a scan status just as s.step does.
//func (s *scanner) eof() int {
//	if s.err != nil {
//		return scanError
//	}
//	if s.endTop {
//		return scanEnd
//	}
//	s.step(s, ' ')
//	if s.endTop {
//		return scanEnd
//	}
//	if s.err == nil {
//		s.err = &SyntaxError{"unexpected end of JSON input", s.bytes}
//	}
//	return scanError
//}
//
//// checkValid verifies that data is valid JSON-encoded data.
//// scan is passed in for use by checkValid to avoid an allocation.
//func checkValid(data []byte, scan *scanner) error {
//	scan.reset()
//	for _, c := range data {
//		scan.bytes++
//		if scan.step(scan, c) == scanError {
//			return scan.err
//		}
//	}
//	if scan.eof() == scanError {
//		return scan.err
//	}
//	return nil
//}
//
//func Unmarshal(data []byte, v interface{}) error {
//	// Check for well-formedness.
//	// Avoids filling out half a data structure
//	// before discovering a JSON syntax error.
//	var d decodeState
//	err := checkValid(data, &d.scan)
//	if err != nil {
//		return err
//	}
//
//	d.init(data)
//	return d.unmarshal(v)
//}

// ReadResultCmd unmarshalls the json reply with data struct for specific
// commands or an interface if it is not a command where we already have a
// struct ready.
func ReadResultCmd(cmd string, message []byte) (flojson.Reply, error) {
	var result flojson.Reply
	//fmt.Println(result) //{<nil> <nil> <nil>}
	var err error
	var objmap map[string]json.RawMessage
	//fmt.Println(objmap) //map[]
	//fmt.Print(message) //[67 108 105 101 110 116 32 115 101 110 116 32 97 110 32 72 84 84 80 32 114 101 113 117 101 115 116 32 116 111 32 97 110 32 72 84 84 80 83 32 115 101 114 118 101 114 46 10]
	err = json.Unmarshal(message, &objmap)
	//fmt.Println(err) //invalid character 'C' looking for beginning of value
	//-1: This implementation does not implement wallet commands
	//fmt.Println(objmap) //map[result:[110 117 108 108] error:[123 34 99 111 100 101 34 58 45 49 44 34 109 101 115 115 97 103 101 34 58 34 84 104 105 115 32 105 109 112 108 101 109 101 110 116 97 116 105 111 110 32 100 111 101 115 32 110 111 116 32 105 109 112 108 101 109 101 110 116 32 119 97 108 108 101 116 32 99 111 109 109 97 110 100 115 34 125] id:[48]]
	if err != nil {
		if strings.Contains(string(message), "401 Unauthorized.") {
			err = fmt.Errorf("Authentication error.")
		} else {
			err = fmt.Errorf("Error unmarshalling json reply: %v", err)
		}
		return result, err
	}
	// Take care of the parts that are the same for all replies.
	var jsonErr flojson.Error
	var id interface{}
	err = json.Unmarshal(objmap["error"], &jsonErr)
	//fmt.Println(jsonErr) //-1: This implementation does not implement wallet commands
	//fmt.Println(err) //<nil>
	if err != nil {
		err = fmt.Errorf("Error unmarshalling json reply: %v", err)
		return result, err
	}
	err = json.Unmarshal(objmap["id"], &id)
	//fmt.Println(id) //0
	//fmt.Println(err) //<nil>
	if err != nil {
		err = fmt.Errorf("Error unmarshalling json reply: %v", err)
		return result, err
	}
	//fmt.Println(cmd) //signmessage

	// If it is a command where we have already worked out the reply,
	// generate put the results in the proper structure.
	// We handle the error condition after the switch statement.
	switch cmd {
	case "createmultisig":
		var res *flojson.CreateMultiSigResult
		err = json.Unmarshal(objmap["result"], &res)
		if err == nil {
			result.Result = res
		}
	case "decodescript":
		var res *flojson.DecodeScriptResult
		err = json.Unmarshal(objmap["result"], &res)
		if err == nil {
			result.Result = res
		}
	case "getaddednodeinfo":
		// getaddednodeinfo can either return a JSON object or a
		// slice of strings depending on the verbose flag.  Choose the
		// right form accordingly.
		if bytes.IndexByte(objmap["result"], '{') > -1 {
			var res []flojson.GetAddedNodeInfoResult
			err = json.Unmarshal(objmap["result"], &res)
			if err == nil {
				result.Result = res
			}
		} else {
			var res []string
			err = json.Unmarshal(objmap["result"], &res)
			if err == nil {
				result.Result = res
			}
		}
	case "getinfo":
		var res *flojson.InfoResult
		err = json.Unmarshal(objmap["result"], &res)
		if err == nil {
			result.Result = res
		}
	case "getblock":
		// getblock can either return a JSON object or a hex-encoded
		// string depending on the verbose flag.  Choose the right form
		// accordingly.
		if bytes.IndexByte(objmap["result"], '{') > -1 {
			var res *flojson.BlockResult
			err = json.Unmarshal(objmap["result"], &res)
			if err == nil {
				result.Result = res
			}
		} else {
			var res string
			err = json.Unmarshal(objmap["result"], &res)
			if err == nil {
				result.Result = res
			}
		}
	case "getblockchaininfo":
		var res *flojson.GetBlockChainInfoResult
		err = json.Unmarshal(objmap["result"], &res)
		if err == nil {
			result.Result = res
		}
	case "getnettotals":
		var res *flojson.GetNetTotalsResult
		err = json.Unmarshal(objmap["result"], &res)
		if err == nil {
			result.Result = res
		}
	case "getnetworkhashps":
		var res int64
		err = json.Unmarshal(objmap["result"], &res)
		if err == nil {
			result.Result = res
		}
	case "getpeerinfo":
		var res []flojson.GetPeerInfoResult
		err = json.Unmarshal(objmap["result"], &res)
		if err == nil {
			result.Result = res
		}
	case "getrawtransaction":
		// getrawtransaction can either return a JSON object or a
		// hex-encoded string depending on the verbose flag.  Choose the
		// right form accordingly.
		if bytes.IndexByte(objmap["result"], '{') > -1 {
			var res *flojson.TxRawResult
			err = json.Unmarshal(objmap["result"], &res)
			if err == nil {
				if len(res.TxComment) < len(res.FloData) {
					res.TxComment = res.FloData
				}
				result.Result = res
			}
		} else {
			var res string
			err = json.Unmarshal(objmap["result"], &res)
			if err == nil {
				result.Result = res
			}
		}
	case "decoderawtransaction":
		var res *flojson.TxRawDecodeResult
		err = json.Unmarshal(objmap["result"], &res)
		if err == nil {
			if len(res.TxComment) < len(res.FloData) {
				res.TxComment = res.FloData
			}
			result.Result = res
		}
	case "getaddressesbyaccount":
		var res []string
		err = json.Unmarshal(objmap["result"], &res)
		if err == nil {
			result.Result = res
		}
	case "getmininginfo":
		var res *flojson.GetMiningInfoResult
		err = json.Unmarshal(objmap["result"], &res)
		if err == nil {
			result.Result = res
		}
	case "getnetworkinfo":
		var res *flojson.GetNetworkInfoResult
		err = json.Unmarshal(objmap["result"], &res)
		if err == nil {
			if res.LocalAddresses == nil {
				res.LocalAddresses = []flojson.LocalAddressesResult{}
			}
			result.Result = res
		}
	case "getrawmempool":
		// getrawmempool can either return a map of JSON objects or
		// an array of strings depending on the verbose flag.  Choose
		// the right form accordingly.
		if bytes.IndexByte(objmap["result"], '{') > -1 {
			var res map[string]flojson.GetRawMempoolResult
			err = json.Unmarshal(objmap["result"], &res)
			if err == nil {
				result.Result = res
			}
		} else {
			var res []string
			err = json.Unmarshal(objmap["result"], &res)
			if err == nil {
				result.Result = res
			}
		}
	case "gettransaction":
		var res *flojson.GetTransactionResult
		err = json.Unmarshal(objmap["result"], &res)
		if err == nil {
			result.Result = res
		}
	case "getwork":
		// getwork can either return a JSON object or a boolean
		// depending on whether or not data was provided.  Choose the
		// right form accordingly.
		if bytes.IndexByte(objmap["result"], '{') > -1 {
			var res *flojson.GetWorkResult
			err = json.Unmarshal(objmap["result"], &res)
			if err == nil {
				result.Result = res
			}
		} else {
			var res bool
			err = json.Unmarshal(objmap["result"], &res)
			if err == nil {
				result.Result = res
			}
		}
	case "validateaddress":
		var res *flojson.ValidateAddressResult
		err = json.Unmarshal(objmap["result"], &res)
		if err == nil {
			result.Result = res
		}
	case "signrawtransaction":
		var res *flojson.SignRawTransactionResult
		err = json.Unmarshal(objmap["result"], &res)
		if err == nil {
			result.Result = res
		}
	case "listaccounts":
		var res map[string]float64
		err = json.Unmarshal(objmap["result"], &res)
		if err == nil {
			result.Result = res
		}
	case "listreceivedbyaccount":
		var res []flojson.ListReceivedByAccountResult
		err = json.Unmarshal(objmap["result"], &res)
		if err == nil {
			result.Result = res
		}
	case "listreceivedbyaddress":
		var res []flojson.ListReceivedByAddressResult
		err = json.Unmarshal(objmap["result"], &res)
		if err == nil {
			result.Result = res
		}
	case "listsinceblock":
		var res *flojson.ListSinceBlockResult
		err = json.Unmarshal(objmap["result"], &res)
		if err == nil {
			if res.Transactions == nil {
				res.Transactions = []flojson.ListTransactionsResult{}
			}
			result.Result = res
		}
	case "listtransactions":
		var res []flojson.ListTransactionsResult
		err = json.Unmarshal(objmap["result"], &res)
		if err == nil {
			result.Result = res
		}
	case "listunspent":
		var res []flojson.ListUnspentResult
		err = json.Unmarshal(objmap["result"], &res)
		if err == nil {
			result.Result = res
		}
	// For commands that return a single item (or no items), we get it with
	// the correct concrete type for free (but treat them separately
	// for clarity).
	case "getblockcount", "getbalance", "getblockhash", "getgenerate",
		"getconnectioncount", "getdifficulty", "gethashespersec",
		"setgenerate", "stop", "settxfee", "getaccount",
		"getnewaddress", "sendtoaddress", "createrawtransaction",
		"sendrawtransaction", "getbestblockhash", "getrawchangeaddress",
		"sendfrom", "sendmany", "addmultisigaddress", "getunconfirmedbalance",
		"getaccountaddress", "estimatefee", "estimatepriority":
		err = json.Unmarshal(message, &result)
	default:
		// None of the standard Bitcoin RPC methods matched.  Try
		// registered custom command reply parsers.
		//fmt.Println(customCmds[cmd]) //{<nil> <nil> }

		if c, ok := customCmds[cmd]; ok && c.replyParser != nil {
			fmt.Println("reach here")
			var res interface{}
			res, err = c.replyParser(objmap["result"])
			if err == nil {
				result.Result = res
			}
		} else {
			// For anything else put it in an interface.  All the
			// data is still there, just a little less convenient
			// to deal with.
			err = json.Unmarshal(message, &result)
			//fmt.Println(result) //{<nil> -1: This implementation does not implement wallet commands 0xc4204ba150}
			//fmt.Println(err) //<nil>
		}
	}
	if err != nil {
		err = fmt.Errorf("Error unmarshalling json reply: %v", err)
		return result, err
	}
	// Only want the error field when there is an actual error to report.
	//fmt.Println(jsonErr.Code) //-1
	if jsonErr.Code != 0 {
		result.Error = &jsonErr
	}
	//fmt.Println(result.Error) //-1: This implementation does not implement wallet commands
	//fmt.Println(&id) //0xc42024e0c0
	result.Id = &id
	//fmt.Println(result) //{<nil> -1: This implementation does not implement wallet commands 0xc42024e0c0}
	return result, err
}

// 因为测试而改写的程序
// 因为测试而改写的程序
// 因为测试而改写的程序

func sendRPC(cmd flojson.Cmd) (flojson.Reply, error) { // 发送RPC调用
	t := 0
	//fmt.Println(user) //user
	//fmt.Println(pass) //pass
	//fmt.Println(server) //localhost:17317
	//fmt.Println(cmd) //&{0 FS267EyRkAEyLNuSaTdbtkgcjjsXEBYhse QmfJxwEBCbfe5SRQpP1T1jaJrCLMuBSwP8FgpeV54rsPLx-FS267EyRkAEyLNuSaTdbtkgcjjsXEBYhse-1604913686}
	for true {
		//reply, err := flojson.RpcSend(user, pass, server, cmd)
		//fmt.Println(reply) //Error reading json message: Error unmarshalling json reply: invalid character 'C' looking for beginning of value
		//fmt.Println(err) //{<nil> <nil> <nil>}

		reply, err := RpcSend(user, pass, server, cmd)
		//fmt.Println(reply) //{<nil> -1: This implementation does not implement wallet commands 0xc4200d2370}
		//fmt.Println(err) //<nil>
		//fmt.Println(reply.Error) //-1: This implementation does not implement wallet commands
		//fmt.Println(reply.Error.Code) //-1

		if err != nil {
			fmt.Println(reply, err)
			return reply, err
		}
		if reply.Error != nil {
			if (reply.Error.Code == -6 && reply.Error.Message == "Insufficient funds") || // 余额不足
				(reply.Error.Code == -4 && strings.HasPrefix(reply.Error.Message, "This transaction requires a transaction fee of at least")) {
				if t > 20 {
					fmt.Println("It's been 10 minutes, perhaps you're really out of funds")
					return reply, reply.Error
				}
				t++
				fmt.Println("Sleeping 30s for a block to re-confirm balance")
				time.Sleep(30 * time.Second)
				continue
			}
			return reply, reply.Error
		}
		return reply, nil
	}
	panic("the above loop didn't return, something terrible has gone wrong")
}

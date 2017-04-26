package utils

import (
	"bytes"
	"errors"
	"fmt"
)

var (
	ctrl              = "\r\n"
	lineBreak         = []byte("|")
	emptyString       = []byte("")
	ctrlLine          = []byte(ctrl)
	newLine           = []byte("\n")
	newCl             = []byte("\r")
	colon             = byte(':')
	colonSlice        = []byte(":")
	endBracket        = byte('}')
	beginBracket      = byte('{')
	endBracketSlice   = []byte("}")
	beginBracketSlice = []byte("{")
	endColonBracket   = []byte("}:")
	beginColonBracket = []byte(":{")
	endTrace          = []byte("End Trace")
)

// MessageParser defines an interface for a message parser which handles
// parsing of a recieved message slice.
type MessageParser interface {
	Parse([]byte) ([]Message, error)
}

// Messages defines a slice of Message structs.
type Messages []Message

// Message defines a struct that details a specific message piece of a data
// recieved.
type Message struct {
	Command []byte   `json:"command"`
	Data    [][]byte `json:"data"`
}

// String returns a stringified version of the giving message.
func (m Message) String() string {
	return fmt.Sprintf("{Command: %+q, Data: %+q}", m.Command, m.Data)
}

// BlockParser defines a package level parser using the blockMessage specification.
/* The block message specification is based on the idea of a simple text based
   protocol which follows specific rules about messages. These rules which are:

   1. All messages must end with a CRTL line ending `\r\n`.

   2. All message must begin with a  opening(`{`) backet and close with a closing(`}`) closing bracket. eg
	    `{A|U|Runner}`.

	 3. If we need to seperate messages within brackets that contain characters like '(',')','{','}' or just a part is
	    kept whole, then using the message exclude block '()' with a message block will make that happen. eg `{A|U|Runner|(U | F || JR (Read | UR))}\r\n`
			Where `(U | F || JR (Read | UR))` is preserved as a single block when parsed into component blocks.

	 4. Multiplex messages can be included together by combining messages wrapped by the brackets with a semicolon(`:`) in between. eg
	 		`{A|U|Runner}:{+SUBS|R|}\r\n`.
*/
var BlockParser blockMessage

// blockMessage defines a struct for the blockParser.
type blockMessage struct{}

// Parse parses the data data coming in and produces a series of Messages
// based on a base pattern.
func (b blockMessage) Parse(msg []byte) ([]Message, error) {
	if len(msg) == 0 {
		return nil, errors.New("Empty Data Received")
	}

	var messages []Message

	blocks, err := b.SplitMultiplex(msg)
	if err != nil {
		return nil, err
	}

	// fmt.Printf("DBLOCKS: %+q\n", blocks)

	for _, block := range blocks {
		if len(block) == 0 {
			continue
		}

		blockParts := b.SplitParts(block)

		var command []byte
		var data [][]byte

		command = blockParts[0]

		if len(blockParts) > 1 {
			data = blockParts[1:len(blockParts)]
		}

		messages = append(messages, Message{
			Command: command,
			Data:    data,
		})
	}

	return messages, nil
}

// SplitParts splits the parts of a message block divided by a vertical bar(`|`)
// symbol. It spits out the parts into their seperate enitites.
func (blockMessage) SplitParts(msg []byte) [][]byte {
	var dataBlocks [][]byte

	var block []byte

	msg = bytes.TrimPrefix(msg, beginBracketSlice)
	msg = bytes.TrimSuffix(msg, endBracketSlice)
	msgLen := len(msg)

	var excludedBlock bool
	var excludedDebt int

	for i := 0; i < msgLen; i++ {
		item := msg[i]

		if item == '(' {
			excludedBlock = true
			excludedDebt++
		}

		if excludedBlock && item == ')' {
			block = append(block, item)
			excludedDebt--

			if excludedDebt <= 0 {
				excludedBlock = false
			}

			continue
		}

		if excludedBlock {
			block = append(block, item)
			continue
		}

		if item == '|' {
			dataBlocks = append(dataBlocks, block)
			block = nil
			continue
		}

		block = append(block, item)
	}

	if len(block) != 0 {
		dataBlocks = append(dataBlocks, block)
	}

	return dataBlocks
}

// SplitMultiplex will split messages down into their separate parts by
// breaking patterns of message packs into their seperate parts.
// multiplex message: `{A|U|Runner}:{+SUBS|R|}\r\n` => []{[]byte("A|U|Runner}\r\n"), []byte("+SUBS|R|bucks\r\n")}.
func (blockMessage) SplitMultiplex(msg []byte) ([][]byte, error) {
	var blocks [][]byte

	var block []byte
	var messageStart bool

	msgLen := len(msg)

	if msg[0] != beginBracket {
		return nil, errors.New("Expected message block/blocks should be enclosed in '{','}'")
	}

	i := -1

	{
	blockLoop:
		for {
			if msg == nil || len(msg) == 0 {
				break blockLoop
			}

			i++

			if i >= len(msg) {
				break blockLoop
			}

			item := msg[i]

			switch item {
			case ':':
				if i+1 >= msgLen {
					return nil, errors.New("blocks must follow {contents} rule")
				}

				piece := msg[i-1]
				var nxtPiece []byte

				if i < msgLen && i+2 < msgLen {
					nxtPiece = msg[i : i+2]
				}

				if -1 != bytes.Compare(nxtPiece, beginColonBracket) && piece == endBracket {
					blocks = append(blocks, block)
					block = nil
					messageStart = false
					continue
				}

				block = append(block, item)
				continue

			case '{':
				// fmt.Printf("BEGINSTART: %+q -> %+q : %+q \n", item, block, msg)
				if messageStart {
					block = append(block, item)
					continue
				}

				messageStart = true
				block = append(block, item)
				continue

			case '}':
				// fmt.Printf("ClOSESTART: %+q -> %+q : %+q \n", item, block, msg)

				if !messageStart {
					return nil, errors.New("Invalid Start of block")
				}

				// Are we at message end and do are we starting a new block?
				if i+1 >= msgLen {
					blocks = append(blocks, block)
					block = nil
					break
				}

				if msg[i+1] == ':' && i+2 >= msgLen {
					return nil, errors.New("Invalid new Block start")
				}

				// fmt.Printf("MSG: %+q -> %+q : %+q\n", msg, block, item)
				// fmt.Printf("MSGD: %+q -> %+q : %+q : %+q\n", msg, msg[i+1:i+2], msg[i+1:i+3], msg[i+3:])

				if msg[i+1] == ':' && msg[i+2] == '{' {
					block = append(block, item)
					blocks = append(blocks, block)
					block = nil
					messageStart = false
					continue
				}

				// We are probably facing a new section without
				if item == '}' && bytes.Equal(msg[i+1:i+3], ctrlLine) {
					block = append(block, item)
					blocks = append(blocks, block)
					block = nil
					messageStart = false

					msg = msg[i+3:]
					i = -1
					// fmt.Printf("MSGOD: %+q -> %+q : %+q \n", item, block, msg)

					if len(msg) == 0 {
						break blockLoop
					}

					continue blockLoop
				}

				ctrl := msg[i+1 : msgLen]
				if 0 == bytes.Compare(ctrl, ctrlLine) {
					if item == endBracket {
						block = append(block, item)
					}

					blocks = append(blocks, block)
					block = nil
					break blockLoop
				}

				block = append(block, item)

			default:

				ctrl := msg[i:msgLen]
				if 0 == bytes.Compare(ctrl, ctrlLine) {
					break blockLoop
				}

				block = append(block, item)
			}
		}

	}

	if len(block) > 1 {
		blocks = append(blocks, block)
	}

	return blocks, nil
}

//==============================================================================

// MakeMessage wraps each byte slice in the multiple byte slice with with
// given header returning a single byte slice joined with a colon : symbol.
func MakeMessage(header string, msgs ...string) []byte {
	var msg []byte

	if msgs != nil {
		var bmsg [][]byte

		for _, msg := range msgs {
			bmsg = append(bmsg, []byte(msg))
		}

		msg = WrapBlockParts(bmsg)
	}

	return WrapCTRLLine(WrapBlock(WrapWithHeader([]byte(header), msg)))
}

// MakeByteMessage wraps each byte slice in the multiple byte slice with with
// given header returning a single byte slice joined with a colon : symbol.
func MakeByteMessage(header []byte, msgs ...[]byte) []byte {
	var msg []byte

	if msgs != nil {
		msg = WrapBlockParts(msgs)
	}

	return WrapCTRLLine(WrapBlock(WrapWithHeader(header, msg)))
}

// MakeMessageGroup joins the provided message strings with a ':' character.
func MakeMessageGroup(msgs ...string) []byte {
	var bmsg [][]byte

	for _, msg := range msgs {
		bmsg = append(bmsg, []byte(msg))
	}

	return WrapCTRLLine(bytes.Join(bmsg, colonSlice))
}

// MakeByteMessageGroup joins series of messages with a ':' character.
func MakeByteMessageGroup(msg ...[]byte) []byte {
	return WrapCTRLLine(bytes.Join(msg, colonSlice))
}

// JoinMessages joins all the provided messages details returning a single byte
// slice.
func JoinMessages(mgs ...Message) []byte {
	var bm [][]byte

	for _, msg := range mgs {
		bm = append(bm, WrapResponseBlock(msg.Command, msg.Data...))
	}

	return bytes.Join(bm, colonSlice)
}

// WrapCTRLLine wraps the giving message with a \r\n ending if not already there.
func WrapCTRLLine(msg []byte) []byte {
	if bytes.HasSuffix(msg, ctrlLine) {
		return msg
	}

	return append(msg, ctrlLine...)
}

// WrapBlock wraps a message in a opening and closing bracket.
func WrapBlock(msg []byte) []byte {
	if bytes.HasPrefix(msg, beginBracketSlice) && bytes.HasSuffix(msg, endBracketSlice) {
		return msg
	}

	return bytes.Join([][]byte{[]byte("{"), msg, []byte("}")}, emptyString)
}

// WrapBlockParts wraps the provided bytes slices with a | symbol.
func WrapBlockParts(msg [][]byte) []byte {
	if len(msg) == 0 {
		return nil
	}

	return bytes.Join(msg, lineBreak)
}

// WrapWithHeader wraps a message with a response wit the given header
// returning a single byte slice.
func WrapWithHeader(header []byte, msg []byte) []byte {
	if len(msg) == 0 {
		return header
	}

	return bytes.Join([][]byte{header, msg}, lineBreak)
}

// WrapResponses wraps each byte slice in the multiple byte slice with with
// given header returning a single byte slice joined with a colon : symbol.
func WrapResponses(header []byte, msgs ...[][]byte) []byte {
	var responses [][]byte

	for _, blocks := range msgs {
		var mod []byte

		if len(header) == 0 {
			mod = WrapWithHeader(header, WrapBlockParts(blocks))
		} else {
			mod = WrapBlockParts(blocks)
		}

		responses = append(
			responses,
			WrapBlock(mod),
		)
	}

	return bytes.Join(responses, colonSlice)
}

// WrapResponseBlock wraps each byte slice in the multiple byte slice with with
// given header returning a single byte slice joined with a colon : symbol.
func WrapResponseBlock(header []byte, msgs ...[]byte) []byte {
	var msg []byte

	if msgs != nil {
		msg = WrapBlockParts(msgs)
	}

	return WrapBlock(WrapWithHeader(header, msg))
}

// WrapResponse wraps each byte slice in the multiple byte slice with with
// given header returning a single byte slice joined with a colon : symbol.
func WrapResponse(header []byte, msgs ...[]byte) []byte {
	var responses [][]byte

	for _, block := range msgs {

		if header != nil {
			block = WrapWithHeader(header, block)
		}

		responses = append(
			responses,
			WrapBlock(block),
		)
	}

	return bytes.Join(responses, colonSlice)
}

// Package localtalk provides basic LocalTalk support.  The package provides
// two layers: Listeners, which are responsible for pulling LLAP packets off
// the "wire", and the Port, which is responsible for decoding LLAP packets,
// dealing with address acquisition, and doing the "logical" stuff.
//
// At the moment, the only listener available is for LocalTalk tunneled over
// UDP, but it would be nice to change that in future.
package localtalk

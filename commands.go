package redimock

import (
	"strings"
)

// TODO : I don't want to fall in `implement another redis` trap. so be careful :)

// ExpectQuit try to return quit command
func (s *Server) ExpectQuit() *Command {
	return s.Expect("Quit").CloseConnection().WillReturn("OK")
}

// ExpectGet return a redis GET command
func (s *Server) ExpectGet(key string, exists bool, result string) *Command {
	c := s.Expect("GET").WithArgs(key)
	if exists {
		return c.WillReturn(BulkString(result))
	}
	return c.WillReturn(nil)
}

// ExpectSet return a redis set command. success could be false only for NX or XX option,
// otherwise it dose not make sense
func (s *Server) ExpectSet(key string, value string, success bool, extra ...string) *Command {
	args := append([]string{key, value}, extra...)
	c := s.Expect("SET").WithArgs(args...)
	if success {
		return c.WillReturn("OK")
	}
	return c.WillReturn(func(args ...string) []interface{} {
		for _, i := range args {
			x := strings.ToUpper(i)
			if x == "NX" || x == "XX" {
				return []interface{}{nil}
			}
		}
		return []interface{}{"OK"}
	})
}

// ExpectPing is the ping command
func (s *Server) ExpectPing() *Command {
	return s.Expect("PING").WithAnyArgs().WillReturn(func(args ...string) []interface{} {
		if len(args) == 0 {
			return []interface{}{"PONG"}
		} else if len(args) == 1 {
			return []interface{}{args[0]}
		}
		return []interface{}{Error("ERR wrong number of arguments for 'ping' command")}
	})
}

// ExpectHSet is the command HSET, if the update is true, then it means the key
// was there already
func (s *Server) ExpectHSet(key, field, value string, update bool) *Command {
	ret := 1
	if update {
		ret = 0
	}
	return s.Expect("HSET").WithArgs(key, field, value).WillReturn(ret)
}

// ExpectHGetAll return the HGETALL command
func (s *Server) ExpectHGetAll(key string, ret map[string]string) *Command {
	arr := make([]BulkString, 0, len(ret)*2)
	for i := range ret {
		arr = append(arr, BulkString(i), BulkString(ret[i]))
	}
	return s.Expect("HGETALL").WithArgs(key).WillReturn(arr)
}

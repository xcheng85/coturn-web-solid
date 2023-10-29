package service

import (
	"fmt"
)

type EmptyExternalIpErr struct {

}

func (r *EmptyExternalIpErr) Error() string {
  return fmt.Sprintf("no external ips of load balancer(s) are available")
}

// assert style in golang
func (s *EmptyExternalIpErr) Is(target error) bool {
  _, ok := target.(*EmptyExternalIpErr)
  if !ok {
    return false
  }
  return true
}

func NewEmptyExternalIpErr() *EmptyExternalIpErr {
	return &EmptyExternalIpErr{}
}
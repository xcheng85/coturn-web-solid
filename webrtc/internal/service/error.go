package service

import (
	"fmt"
)

type InvalidExternalIpErr struct {
  address string
}

func (r *InvalidExternalIpErr) Error() string {
  return fmt.Sprintf("External Ip: %s is invalid", r.address)
}

// assert style in golang
func (s *InvalidExternalIpErr) Is(target error) bool {
  _, ok := target.(*InvalidExternalIpErr)
  if !ok {
    return false
  }
  return true
}

func NewInvalidExternalIpErr(address string) *InvalidExternalIpErr {
	return &InvalidExternalIpErr{address}
}
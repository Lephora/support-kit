package unit_test

import (
	. "gopkg.in/check.v1"
	"testing"
	"thoughtworks.com/lephora/support-kit/accecptance-test/test_assertion"
)

func Test(t *testing.T) { TestingT(t) }

type CheckerSuite struct{}

var _ = Suite(&CheckerSuite{})

func (s *CheckerSuite) TestGenericAssertorEqual(c *C) {
	// given
	assertor := test_assertion.GenericAssertor{}

	// when
	assert, err := assertor.Assert(1, test_assertion.Equals, 1)

	// then
	c.Check(err, IsNil)
	c.Check(assert, Equals, true)

	// when
	assert, err = assertor.Assert("1", test_assertion.Equals, 1)

	// then
	c.Check(err, NotNil)
	c.Check(assert, Equals, false)

	// when
	assert, err = assertor.Assert("yes", test_assertion.Equals, "no")

	// then
	c.Check(err, NotNil)
	c.Check(assert, Equals, false)

	// when
	assert, err = assertor.Assert("yes", test_assertion.Equals, "yes")

	// then
	c.Check(err, IsNil)
	c.Check(assert, Equals, true)
}

func (s *CheckerSuite) TestGenericAssertorMatch_01(c *C) {
	// given
	assertor := test_assertion.GenericAssertor{}

	// when
	assert, err := assertor.Assert("1", test_assertion.Match, "1")

	// then
	c.Check(err, IsNil)
	c.Check(assert, Equals, true)
}

func (s *CheckerSuite) TestGenericAssertorMatch_02(c *C) {
	// given
	assertor := test_assertion.GenericAssertor{}

	// when
	assert, err := assertor.Assert("hello, bob", test_assertion.Match, "hello, bob")

	// then
	c.Check(err, IsNil)
	c.Check(assert, Equals, true)
}

func (s *CheckerSuite) TestGenericAssertorMatch_03(c *C) {
	// given
	assertor := test_assertion.GenericAssertor{}

	// when
	assert, err := assertor.Assert("hello, bob", test_assertion.Match, "^hello, bob$")

	// then
	c.Check(err, IsNil)
	c.Check(assert, Equals, true)
}

func (s *CheckerSuite) TestGenericAssertorMatch_04(c *C) {
	// given
	assertor := test_assertion.GenericAssertor{}

	// when
	assert, err := assertor.Assert("hello, bob", test_assertion.Match, `^[a-zA-Z,\s]+$`)

	// then
	c.Check(err, IsNil)
	c.Check(assert, Equals, true)
}

func (s *CheckerSuite) TestGenericAssertorMatch_05(c *C) {
	// given
	assertor := test_assertion.GenericAssertor{}

	// when
	assert, err := assertor.Assert("hello, bob", test_assertion.Match, `^[a-zA-Z,]+$`)

	// then
	c.Check(err, IsNil)
	c.Check(assert, Equals, false)
}

func (s *CheckerSuite) TestGenericAssertorMatch_06(c *C) {
	// given
	assertor := test_assertion.GenericAssertor{}

	// when
	assert, err := assertor.Assert("10086", test_assertion.Match, `^[0-9]+$`)

	// then
	c.Check(err, IsNil)
	c.Check(assert, Equals, true)
}

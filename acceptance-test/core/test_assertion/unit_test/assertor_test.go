package unit

import (
	. "gopkg.in/check.v1"
	"thoughtworks.com/lephora/support-kit/accecptance-test/test_assertion"
)

type AssertorSuite struct{}

var _ = Suite(&AssertorSuite{})

func (s *AssertorSuite) TestJsonAssertor(c *C) {
	// given
	assertor := test_assertion.JsonAssertor{}
	object := `
{
  "name": {"first": "Tom", "last": "Anderson"},
  "age":37,
  "children": ["Sara","Alex","Jack"],
  "fav.movie": "Deer Hunter",
  "friends": [
    {"first": "Dale", "last": "Murphy", "age": 44, "nets": ["ig", "fb", "tw"]},
    {"first": "Roger", "last": "Craig", "age": 68, "nets": ["fb", "tw"]},
    {"first": "Jane", "last": "Murphy", "age": 47, "nets": ["ig", "tw"]}
  ]
}
`
	// when
	assert, err := assertor.Object(object).Assert("name.last", test_assertion.Equals, "Anderson")
	// then
	c.Check(err, IsNil)
	c.Check(assert, Equals, true)

	// when
	assert, err = assertor.Object(object).Assert("age", test_assertion.Equals, "37")
	// then
	c.Check(err, IsNil)
	c.Check(assert, Equals, true)

	// when
	assert, err = assertor.Object(object).Assert("children", test_assertion.Equals, `["Sara","Alex","Jack"]`)
	// then
	c.Check(err, IsNil)
	c.Check(assert, Equals, true)

	// when
	assert, err = assertor.Object(object).Assert("friends.#.first", test_assertion.Equals, `["Dale","Roger","Jane"]`)
	// then
	c.Check(err, IsNil)
	c.Check(assert, Equals, true)

	// when
	assert, err = assertor.Object(object).Assert("friends.1.last", test_assertion.Equals, `Craig`)
	// then
	c.Check(err, IsNil)
	c.Check(assert, Equals, true)
}

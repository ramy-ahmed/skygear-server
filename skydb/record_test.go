package skydb

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRecord(t *testing.T) {
	Convey("Set transient field", t, func() {
		note0 := Record{
			ID: NewRecordID("note", "0"),
			Transient: Data{
				"content": "hello world",
			},
		}

		So(note0.Get("content"), ShouldBeNil)
		So(note0.Get("_transient"), ShouldResemble, Data{
			"content": "hello world",
		})
		So(note0.Get("_transient_content"), ShouldEqual, "hello world")
	})

	Convey("Set transient field", t, func() {
		note0 := Record{
			ID: NewRecordID("note", "0"),
		}

		note0.Set("_transient", Data{
			"content": "hello world",
		})

		So(note0.Data["content"], ShouldBeNil)
		So(note0.Transient, ShouldResemble, Data{
			"content": "hello world",
		})
	})

	Convey("Set individual transient field", t, func() {
		note0 := Record{
			ID: NewRecordID("note", "0"),
			Transient: Data{
				"existing": "should be here",
			},
		}

		note0.Set("_transient_content", "hello world")

		So(note0.Data["content"], ShouldBeNil)
		So(note0.Transient, ShouldResemble, Data{
			"content":  "hello world",
			"existing": "should be here",
		})
	})
}

func TestRecordACL(t *testing.T) {
	Convey("Record with ACL", t, func() {
		userinfo := &UserInfo{
			ID:    "user1",
			Roles: []string{"admin"},
		}

		stranger := &UserInfo{
			ID:    "stranger",
			Roles: []string{"nobody"},
		}

		Convey("Check access right base on role", func() {
			note := Record{
				ID:         NewRecordID("note", "0"),
				DatabaseID: "",
				ACL: RecordACL{
					NewRecordACLEntryRole("admin", ReadLevel),
				},
			}

			So(note.Accessible(userinfo, ReadLevel), ShouldBeTrue)
			So(note.Accessible(stranger, ReadLevel), ShouldBeFalse)
		})

		Convey("Check access right base on direct ace", func() {
			note := Record{
				ID:         NewRecordID("note", "0"),
				DatabaseID: "",
				ACL: RecordACL{
					NewRecordACLEntryDirect("user1", ReadLevel),
				},
			}

			So(note.Accessible(userinfo, ReadLevel), ShouldBeTrue)
			So(note.Accessible(stranger, ReadLevel), ShouldBeFalse)
		})

		Convey("Grant permission on any ACE matched", func() {
			note := Record{
				ID:         NewRecordID("note", "0"),
				DatabaseID: "",
				ACL: RecordACL{
					NewRecordACLEntryDirect("stranger", ReadLevel),
					NewRecordACLEntryRole("admin", ReadLevel),
				},
			}

			So(note.Accessible(userinfo, ReadLevel), ShouldBeTrue)
			So(note.Accessible(stranger, ReadLevel), ShouldBeTrue)
		})

		Convey("Write permission superset read permission", func() {
			note := Record{
				ID:         NewRecordID("note", "0"),
				DatabaseID: "",
				ACL: RecordACL{
					NewRecordACLEntryDirect("stranger", WriteLevel),
					NewRecordACLEntryRole("admin", WriteLevel),
				},
			}
			So(note.Accessible(userinfo, ReadLevel), ShouldBeTrue)
			So(note.Accessible(stranger, ReadLevel), ShouldBeTrue)
		})

		Convey("Reject write on read only permission", func() {
			note := Record{
				ID:         NewRecordID("note", "0"),
				DatabaseID: "",
				ACL: RecordACL{
					NewRecordACLEntryDirect("stranger", ReadLevel),
					NewRecordACLEntryRole("admin", ReadLevel),
				},
			}

			So(note.Accessible(userinfo, WriteLevel), ShouldBeFalse)
			So(note.Accessible(stranger, WriteLevel), ShouldBeFalse)
		})
	})
}

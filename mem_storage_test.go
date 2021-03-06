package main

import (
	proto "github.com/huin/mqtt"
	"reflect"
	"testing"
)

func TestUpdateRetain(t *testing.T) {
	//var topic string
	topic := "/a/b/c"
	mem := NewMemStorage()

	origm := &proto.Publish{
		Header: proto.Header{
			DupFlag:  false,
			QosLevel: proto.QosAtMostOnce,
			Retain:   false,
		},
		TopicName: "/a/b/c",
		Payload:   proto.BytesPayload{1, 2, 3},
	}
	mem.UpdateRetain(topic, origm)
	if m, err := mem.GetRetain(topic); err == nil {
		if origm.TopicName != m.TopicName {
			t.Errorf("could not update: %v", topic)
		}
	} else {
		t.Errorf("could not update: %v", topic)
	}

	// Same topic Again
	origm.MessageId = 0x1234
	mem.UpdateRetain(topic, origm)
	if m, err := mem.GetRetain(topic); err == nil {
		if m.MessageId != 0x1234 {
			t.Errorf("could not update: %v", topic)
		}
	} else {
		t.Errorf("could not update: %v", topic)
	}

	// Other topic
	topic2 := "/other/topic"
	m2 := &proto.Publish{
		Header: proto.Header{
			DupFlag:  false,
			QosLevel: proto.QosAtMostOnce,
			Retain:   false,
		},
		TopicName: topic2,
		Payload:   proto.BytesPayload{1, 2, 3},
	}
	mem.UpdateRetain(topic2, m2)
	if m, err := mem.GetRetain(topic2); err == nil {
		if m.TopicName != topic2 {
			t.Errorf("could not update: %v", topic2)
		}
		// check existing topic
		if mm, err := mem.GetRetain(topic); err == nil {
			if mm.TopicName == m.TopicName {
				t.Errorf("Duplicated message")
			}
		} else {
			t.Errorf("could not update after other topic: %v", topic)
		}

	} else {
		t.Errorf("could not update: %v", topic2)
	}

}

func TestSubscribe(t *testing.T) {
	var topic string
	var cid string
	var cid2 string
	mem := NewMemStorage()

	topic = "/a/b/c"

	cid = "cid1"
	mem.Subscribe(topic, cid)
	if reflect.DeepEqual(mem.TopicTable[topic], []string{cid}) != true {
		t.Errorf("clientid not found: %v", cid)
	}

	cid2 = "cid2"
	mem.Subscribe(topic, cid2)
	if reflect.DeepEqual(mem.TopicTable[topic], []string{cid, cid2}) != true {
		t.Errorf("clientid not found: %v", cid2)
	}

	mem.Unsubscribe(topic, cid)
	for _, c := range mem.TopicTable[topic] {
		if c == cid {
			t.Errorf("Unsubscribed id found: %v", cid)
		}
	}

	mem.Unsubscribe(topic, cid2)
	for _, c := range mem.TopicTable[topic] {
		if c == cid2 {
			t.Errorf("Unsubscribed id found: %v", cid2)
		}
	}

}
func TestcreateStoredMsgId(t *testing.T) {
	m := &proto.Publish{
		Header: proto.Header{
			DupFlag:  false,
			QosLevel: proto.QosAtMostOnce,
			Retain:   false,
		},
		TopicName: "/a/b/c",
		Payload:   proto.BytesPayload{1, 2, 3},
	}
	if createStoredMsgId("a", m) != "a-0" {
		t.Errorf("StoredMsgId creating failed with nomsg")
	}

	var msgid uint16
	msgid = 1000
	m2 := &proto.Publish{
		Header: proto.Header{
			DupFlag:  false,
			QosLevel: proto.QosAtMostOnce,
			Retain:   false,
		},
		MessageId: msgid,
		TopicName: "/a/b/c",
		Payload:   proto.BytesPayload{1, 2, 3},
	}
	if createStoredMsgId("a", m2) != "a-1000" {
		t.Errorf("StoredMsgId creating failed")
	}

}

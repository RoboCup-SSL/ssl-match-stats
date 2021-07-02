package matchstats

import (
	"github.com/pkg/errors"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"os"
)

type Collector struct {
	Collection *MatchStatsCollection
}

func NewCollector() *Collector {
	generator := new(Collector)
	generator.Collection = new(MatchStatsCollection)
	generator.Collection.MatchStats = []*MatchStats{}
	return generator
}

func (a *Collector) Process(filename string) error {
	generator := NewGenerator()

	matchStats, err := generator.Process(filename)
	if err != nil {
		return errors.Wrap(err, "Could not create match states")
	} else {
		a.Collection.MatchStats = append(a.Collection.MatchStats, matchStats)
	}
	return nil
}

func (a *Collector) WriteJson(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return errors.Wrap(err, "Could not create JSON output file")
	}

	jsonMarsh := protojson.MarshalOptions{EmitUnpopulated: true, Indent: "  "}
	if data, err := jsonMarsh.Marshal(a.Collection); err != nil {
		return errors.Wrap(err, "Could not marshal match stats to json")
	} else if _, err := f.Write(data); err != nil {
		return errors.Wrap(err, "Could write marshaled data to file")
	}
	return f.Close()
}

func (a *Collector) WriteBin(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return errors.Wrap(err, "Could not create Binary output file")
	}

	bytes, err := proto.Marshal(a.Collection)
	if err != nil {
		return errors.Wrap(err, "Could not marshal match stats to binary")
	}
	_, err = f.Write(bytes)
	if err != nil {
		return errors.Wrap(err, "Could not write match stats to binary")
	}
	return f.Close()
}

func (a *Collector) ReadBin(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	err = proto.Unmarshal(bytes, a.Collection)
	if err != nil {
		return err
	}

	return nil
}

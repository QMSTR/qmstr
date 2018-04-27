package config

import (
	"testing"
)

func TestCompleteConfig(t *testing.T) {

	var config = `
package:
  name: "The Test"
  metadata:
    Vendor: "Endocode"
    OcFossLiaison: "Mirko Boehm"
    OcComplianceContact: "foss@endocode.com"
  server:
    rpcaddress: ":12345"
    dbaddress: "testhost:54321"
    dbworkers: 4
  analysis:
    - analyzer: test-analyzer
      name: "The Testalyzer"
      selector: sourcecode
      pathsub:
        - old: "/the/path"
          new: "/buildroot"
      config:
        workdir: "/buildroot"
    - analyzer: test-analyzer-2
      name: "The Testalyzer 2"
      selector: sourcecode
      pathsub:
        - old: "/the/path"
          new: "/buildroot"
      config:
        workdir: "/buildroot"
        testfile: "/the/test"
  reporting:
    - reporter: test-reporter
      name: "The test reporter"
      config:
        tester: "Endocode"
`
	_, err := readConfig([]byte(config))
	if err != nil {
		t.Logf("Broken config %v", err)
		t.Fail()
	}
}

func TestMissingModuleInstanceName(t *testing.T) {

	var config = `
package:
  name: "The Test"
  metadata:
    Vendor: "Endocode"
    OcFossLiaison: "Mirko Boehm"
    OcComplianceContact: "foss@endocode.com"
  server:
    rpcaddress: ":12345"
    dbaddress: "testhost:54321"
    dbworkers: 4
  analysis:
    - analyzer: test-analyzer
      name: "The Testalyzer"
      selector: sourcecode
      pathsub:
        - old: "/the/path"
          new: "/buildroot"
      config:
        workdir: "/buildroot"
    - analyzer: test-analyzer-2
      name: "The Testalyzer 2"
      selector: sourcecode
      pathsub:
        - old: "/the/path"
          new: "/buildroot"
      config:
        workdir: "/buildroot"
        testfile: "/the/test"
  reporting:
    - reporter: test-reporter
      config:
        tester: "Endocode"
`
	_, err := readConfig([]byte(config))
	if err == nil || err.Error() != "1. reporter misconfigured Name invalid" {
		t.Log(err)
		t.Fail()
	}
}

func TestDuplicateModuleInstanceNames(t *testing.T) {

	var config = `
package:
  name: "The Test"
  metadata:
    Vendor: "Endocode"
    OcFossLiaison: "Mirko Boehm"
    OcComplianceContact: "foss@endocode.com"
  server:
    rpcaddress: ":12345"
    dbaddress: "testhost:54321"
    dbworkers: 4
  analysis:
    - analyzer: test-analyzer
      name: "The Testalyzer"
      selector: sourcecode
      pathsub:
        - old: "/the/path"
          new: "/buildroot"
      config:
        workdir: "/buildroot"
    - analyzer: test-analyzer-2
      name: "The Testalyzer"
      selector: sourcecode
      pathsub:
        - old: "/the/path"
          new: "/buildroot"
      config:
        workdir: "/buildroot"
        testfile: "/the/test"
  reporting:
    - reporter: test-reporter
      config:
        tester: "Endocode"
`
	_, err := readConfig([]byte(config))
	if err == nil || err.Error() != "2. analyzer misconfigured duplicate value of The Testalyzer in Name" {
		t.Log(err)
		t.Fail()
	}
}

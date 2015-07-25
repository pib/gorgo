package vm_test

import (
	"bytes"
	"path/filepath"

	. "github.com/pib/gorgo/vm"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("VM", func() {
	var (
		modPath string
		vm      *VM
		out     *bytes.Buffer
		err     error
	)
	JustBeforeEach(func() {
		vm, err = NewVM(modPath, false)
		Expect(err).NotTo(HaveOccurred())
		out = new(bytes.Buffer)
		vm.Stdout = out
		err = vm.Run()
	})

	Describe("Running myfuncs.pyc", func() {
		BeforeEach(func() {
			modPath, err = filepath.Abs("../cmd/gorgo/testapps/myfuncs.pyc")
			Expect(err).NotTo(HaveOccurred())
		})
		It("Generates the expected output", func() {
			expected := "Hello Go ! \n5^2 = 25 \n"
			Expect(out.String()).To(Equal(expected))
		})
	})
})

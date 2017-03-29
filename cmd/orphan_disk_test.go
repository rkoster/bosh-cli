package cmd_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/cloudfoundry/bosh-cli/cmd"
	fakedir "github.com/cloudfoundry/bosh-cli/director/directorfakes"
	fakeui "github.com/cloudfoundry/bosh-cli/ui/fakes"
)

var _ = Describe("OrphanDiskCmd", func() {
	var (
		ui       *fakeui.FakeUI
		director *fakedir.FakeDirector
		command  OrphanDiskCmd
	)

	BeforeEach(func() {
		ui = &fakeui.FakeUI{}
		director = &fakedir.FakeDirector{}
		command = NewOrphanDiskCmd(ui, director)
	})

	Describe("Run", func() {
		var (
			opts OrphanDiskOpts
		)

		BeforeEach(func() {
			opts = OrphanDiskOpts{
				Args: OrphanDiskArgs{CID: "disk-cid"},
			}
		})

		act := func() error { return command.Run(opts) }

		It("orphans disk", func() {
			disk := &fakedir.FakeOrphanedDisk{}
			director.FindOrphanedDiskReturns(disk, nil)

			err := act()
			Expect(err).ToNot(HaveOccurred())

			Expect(director.FindOrphanedDiskArgsForCall(0)).To(Equal("disk-cid"))
			Expect(disk.OrphanCallCount()).To(Equal(1))
		})

		It("returns error if orphaning disk failed", func() {
			disk := &fakedir.FakeOrphanedDisk{}
			director.FindOrphanedDiskReturns(disk, nil)

			disk.OrphanReturns(errors.New("fake-err"))

			err := act()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-err"))
		})

		It("does not orphan disk if confirmation is rejected", func() {
			disk := &fakedir.FakeOrphanedDisk{}
			director.FindOrphanedDiskReturns(disk, nil)

			ui.AskedConfirmationErr = errors.New("stop")

			err := act()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("stop"))

			Expect(disk.OrphanCallCount()).To(Equal(0))
		})

		It("returns error if finding disk failed", func() {
			director.FindOrphanedDiskReturns(nil, errors.New("fake-err"))

			err := act()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-err"))
		})
	})
})

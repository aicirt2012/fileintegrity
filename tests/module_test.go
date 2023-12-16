package tests

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/aicirt2012/fileintegrity"
	"github.com/aicirt2012/fileintegrity/tests/common"
)

func TestUpsertFlow(t *testing.T) {
	dir, files := common.CreateScenario("upsert", common.Files{
		common.NewFile(`a\a1.txt`, `2022-05-06T00:40:21+02:00`, `a1 sample txt`),
		common.NewFile(`a\a2.txt`, `2022-05-06T00:40:21+02:00`, `a2 sample txt`),
		common.NewFile(`b\b1.md`, `2022-05-06T00:40:21+02:00`, `b1 sample md`),
		common.NewFile(`b\b2.md`, `2022-05-06T00:40:21+02:00`, `b2 sample md`),
		common.NewFile(`b\bb\bb1.md`, `2022-05-06T00:40:21+02:00`, `bb1 sample md`),
	})

	fileintegrity.Upsert(dir, fileintegrity.EnabledOptions())

	common.AssertFilesExist(t, dir, files)
	common.AssertIntegrityFile(t, dir, []common.FileHash{
		common.NewFileHash(`85b883632e34f9f140915c79f1f4d131f50784a1077c0b1516e38fe226f72301`, ``, `2022-05-06T00:40:21+02:00`, `13`, `a\a1.txt`),
		common.NewFileHash(`2592c50e3d57402c5b5f2293bb2a52dfb38bfc91ae1c9a1f2452b798d53bf7c6`, ``, `2022-05-06T00:40:21+02:00`, `13`, `a\a2.txt`),
		common.NewFileHash(`d64783f26f53c1e668cc75b30f29a89b42e0d19ddddb93bffa1fce509a139922`, ``, `2022-05-06T00:40:21+02:00`, `12`, `b\b1.md`),
		common.NewFileHash(`f59d71f706fb095ce60a3babaf1e5cd65521154ab6854b31d0eb93e678ecddc6`, ``, `2022-05-06T00:40:21+02:00`, `12`, `b\b2.md`),
		common.NewFileHash(`ff6464b4321e5d9b09ae7cb7ba219cee688099f232ef5b978be5f7c94083cc4b`, ``, `2022-05-06T00:40:21+02:00`, `13`, `b\bb\bb1.md`),
	})

	common.AssertUpsertLogFile(t, dir, 0, 5, 0, 0)

	time.Sleep(time.Second)
	common.RemoveFile(dir, `b\bb\bb1.md`)
	fileintegrity.Upsert(dir, fileintegrity.EnabledOptions())

	common.AssertIntegrityFile(t, dir, []common.FileHash{
		common.NewFileHash(`85b883632e34f9f140915c79f1f4d131f50784a1077c0b1516e38fe226f72301`, ``, `2022-05-06T00:40:21+02:00`, `13`, `a\a1.txt`),
		common.NewFileHash(`2592c50e3d57402c5b5f2293bb2a52dfb38bfc91ae1c9a1f2452b798d53bf7c6`, ``, `2022-05-06T00:40:21+02:00`, `13`, `a\a2.txt`),
		common.NewFileHash(`d64783f26f53c1e668cc75b30f29a89b42e0d19ddddb93bffa1fce509a139922`, ``, `2022-05-06T00:40:21+02:00`, `12`, `b\b1.md`),
		common.NewFileHash(`f59d71f706fb095ce60a3babaf1e5cd65521154ab6854b31d0eb93e678ecddc6`, ``, `2022-05-06T00:40:21+02:00`, `12`, `b\b2.md`),
	})
	common.AssertUpsertLogFile(t, dir, 4, 0, 0, 1)

	time.Sleep(time.Second)
	common.UpdateFile(dir, `b\b2.md`, `changed content`, `2023-05-06T13:00:21+02:00`)
	fileintegrity.Upsert(dir, fileintegrity.EnabledOptions())

	common.AssertIntegrityFile(t, dir, []common.FileHash{
		common.NewFileHash(`85b883632e34f9f140915c79f1f4d131f50784a1077c0b1516e38fe226f72301`, ``, `2022-05-06T00:40:21+02:00`, `13`, `a\a1.txt`),
		common.NewFileHash(`2592c50e3d57402c5b5f2293bb2a52dfb38bfc91ae1c9a1f2452b798d53bf7c6`, ``, `2022-05-06T00:40:21+02:00`, `13`, `a\a2.txt`),
		common.NewFileHash(`d64783f26f53c1e668cc75b30f29a89b42e0d19ddddb93bffa1fce509a139922`, ``, `2022-05-06T00:40:21+02:00`, `12`, `b\b1.md`),
		common.NewFileHash(`b92d13bbe02db7ca7686a8e7b854de49c7455948c05cf91a47044278395e212e`, ``, `2023-05-06T13:00:21+02:00`, `15`, `b\b2.md`),
	})
	common.AssertUpsertLogFile(t, dir, 3, 0, 1, 0)
}

func TestUpsertFlow_disabledLog(t *testing.T) {
	dir, files := common.CreateScenario("upsert", common.Files{
		common.NewFile(`a\a1.txt`, `2022-05-06T00:40:21+02:00`, `a1 sample txt`),
	})

	fileintegrity.Upsert(dir, fileintegrity.DisabledOptions())

	common.AssertFilesExist(t, dir, files)
	common.AssertIntegrityFile(t, dir, []common.FileHash{
		common.NewFileHash(`85b883632e34f9f140915c79f1f4d131f50784a1077c0b1516e38fe226f72301`, ``, `2022-05-06T00:40:21+02:00`, `13`, `a\a1.txt`),
	})
	common.AssertLogFileNotExists(t, dir)
}

func TestVerifyFlow_happyCase(t *testing.T) {
	dir, files := common.CreateScenario("verify.happyCase", common.Files{
		common.NewFile(`a\a1.txt`, `2022-05-06T00:40:21+02:00`, `a1 sample txt`),
		common.NewFile(`a\a2.txt`, `2022-05-06T00:40:21+02:00`, `a2 sample txt`),
		common.NewFile(`b\b1.md`, `2022-05-06T00:40:21+02:00`, `b1 sample md`),
		common.NewFile(`b\b2.md`, `2022-05-06T00:40:21+02:00`, `b2 sample md`),
		common.NewFile(`b\bb\bb1.md`, `2022-05-06T00:40:21+02:00`, `bb1 sample md`),
	})

	common.CreateIntegrityFile(t, dir, []common.FileHash{
		common.NewFileHash(`85b883632e34f9f140915c79f1f4d131f50784a1077c0b1516e38fe226f72301`, `2023-05-06T15:12:00.8247784+02:00`, `2022-05-06T00:40:21+02:00`, `13`, `a\a1.txt`),
		common.NewFileHash(`2592c50e3d57402c5b5f2293bb2a52dfb38bfc91ae1c9a1f2452b798d53bf7c6`, `2023-05-06T15:12:00.8230798+02:00`, `2022-05-06T00:40:21+02:00`, `13`, `a\a2.txt`),
		common.NewFileHash(`d64783f26f53c1e668cc75b30f29a89b42e0d19ddddb93bffa1fce509a139922`, `2023-05-06T15:12:00.8242669+02:00`, `2022-05-06T00:40:21+02:00`, `12`, `b\b1.md`),
		common.NewFileHash(`f59d71f706fb095ce60a3babaf1e5cd65521154ab6854b31d0eb93e678ecddc6`, `2023-05-06T15:12:00.8247784+02:00`, `2022-05-06T00:40:21+02:00`, `12`, `b\b2.md`),
		common.NewFileHash(`ff6464b4321e5d9b09ae7cb7ba219cee688099f232ef5b978be5f7c94083cc4b`, `2023-05-06T15:12:00.8247784+02:00`, `2022-05-06T00:40:21+02:00`, `13`, `b\bb\bb1.md`),
	})

	fileintegrity.Verify(dir, fileintegrity.EnabledOptions())

	common.AssertFilesExist(t, dir, files)
	common.AssertVerifyLogFile(t, dir, 5, 0)
}

func TestVerifyFlow_fileNotExists(t *testing.T) {
	dir, files := common.CreateScenario("verify.fileNotExists", common.Files{})

	common.CreateIntegrityFile(t, dir, []common.FileHash{
		common.NewFileHash(`a49b137c40dab92ce8fda591f22f3f5f27b94d750861f68d0558971b00ad33a2`, `2023-05-06T15:12:00.8247784+02:00`, `2022-05-06T00:40:21+02:00`, `12`, `b\bb\bb1.md`),
	})

	fileintegrity.Verify(dir, fileintegrity.EnabledOptions())

	common.AssertFilesExist(t, dir, files)
	common.AssertVerifyLogFile(t, dir, 0, 1)
}

func TestVerifyFlow_fileSizeDiffers(t *testing.T) {
	dir, files := common.CreateScenario("verify.fileSizeDiffers", common.Files{
		common.NewFile(`a\a1.txt`, `2022-05-06T00:40:21+02:00`, `a1 sample txt+`),
	})

	common.CreateIntegrityFile(t, dir, []common.FileHash{
		common.NewFileHash(`a49b137c40dab92ce8fda591f22f3f5f27b94d750861f68d0558971b00ad33a2`, `2023-05-06T15:12:00.8247784+02:00`, `2022-05-06T00:40:21+02:00`, `13`, `a\a1.txt`),
	})

	fileintegrity.Verify(dir, fileintegrity.EnabledOptions())

	common.AssertFilesExist(t, dir, files)
	common.AssertVerifyLogFile(t, dir, 0, 1)
}

func TestVerifyFlow_invalidHash(t *testing.T) {
	dir, files := common.CreateScenario("verify.fileNotExists", common.Files{
		common.NewFile(`a\a1.txt`, `2022-05-06T00:40:21+02:00`, `a1 sample txt`),
	})

	common.CreateIntegrityFile(t, dir, []common.FileHash{
		common.NewFileHash(`85b883632e34f9f140915c79f1f4d131f50784a1077c0b1516e38fe226f72302`, `2023-05-06T15:12:00.8247784+02:00`, `2022-05-06T00:40:21+02:00`, `13`, `a\a1.txt`),
	})

	fileintegrity.Verify(dir, fileintegrity.EnabledOptions())

	common.AssertFilesExist(t, dir, files)
	common.AssertVerifyLogFile(t, dir, 0, 1)
}

func TestVerifyFlow_disabledLog(t *testing.T) {
	dir, files := common.CreateScenario("verify.fileNotExistsNoLogs", common.Files{})

	common.CreateIntegrityFile(t, dir, []common.FileHash{
		common.NewFileHash(`a49b137c40dab92ce8fda591f22f3f5f27b94d750861f68d0558971b00ad33a2`, `2023-05-06T15:12:00.8247784+02:00`, `2022-05-06T00:40:21+02:00`, `12`, `b\bb\bb1.md`),
	})

	fileintegrity.Verify(dir, fileintegrity.DisabledOptions())

	common.AssertFilesExist(t, dir, files)
	common.AssertLogFileNotExists(t, dir)
}

func TestDuplicateFlow(t *testing.T) {
	dir, files := common.CreateScenario("check-duplicates", common.Files{})

	common.CreateIntegrityFile(t, dir, []common.FileHash{
		common.NewFileHash(`85b883632e34f9f140915c79f1f4d131f50784a1077c0b1516e38fe226f72302`, `2023-05-06T15:12:00.8247784+02:00`, `2022-05-06T00:40:21+02:00`, `120`, `unique.md`),
		common.NewFileHash(`a49b137c40dab92ce8fda591f22f3f5f27b94d750861f68d0558971b00ad33a2`, `2023-05-06T15:12:00.8247784+02:00`, `2022-05-06T00:40:21+02:00`, `120`, `duplicate a.md`),
		common.NewFileHash(`a49b137c40dab92ce8fda591f22f3f5f27b94d750861f68d0558971b00ad33a2`, `2023-05-06T15:12:00.8247784+02:00`, `2022-05-06T00:40:21+02:00`, `120`, `duplicate b.md`),
		common.NewFileHash(`a49b137c40dab92ce8fda591f22f3f5f27b94d750861f68d0558971b00ad33a2`, `2023-05-06T15:12:00.8247784+02:00`, `2022-05-06T00:40:21+02:00`, `120`, `duplicate c.md`),
		common.NewFileHash(`ff6464b4321e5d9b09ae7cb7ba219cee688099f232ef5b978be5f7c94083cc4b`, `2023-05-06T15:12:00.8247784+02:00`, `2022-05-06T00:40:21+02:00`, `120`, `duplicate I.md`),
		common.NewFileHash(`ff6464b4321e5d9b09ae7cb7ba219cee688099f232ef5b978be5f7c94083cc4b`, `2023-05-06T15:12:00.8247784+02:00`, `2022-05-06T00:40:21+02:00`, `120`, `duplicate II.md`),
	})

	fileintegrity.CheckDuplicates(dir, fileintegrity.EnabledOptions())

	common.AssertFilesExist(t, dir, files)
	common.AssertDuplicateLogFile(t, dir, []common.LogBlock{
		common.NewDuplicateLogBlock(`a49b137c40dab92ce8fda591f22f3f5f27b94d750861f68d0558971b00ad33a2`, []string{
			`duplicate a.md`,
			`duplicate b.md`,
			`duplicate c.md`,
		}),
		common.NewDuplicateLogBlock(`ff6464b4321e5d9b09ae7cb7ba219cee688099f232ef5b978be5f7c94083cc4b`, []string{
			`duplicate I.md`,
			`duplicate II.md`,
		}),
	}, 3, 3)
}

func TestContainedFlow(t *testing.T) {
	dir, files := common.CreateScenario("check-contained", common.Files{
		common.NewFile(`base\a.txt`, `2022-05-06T00:40:21+02:00`, `unique text1 `+common.StaticContent(101)),
		common.NewFile(`base\b.txt`, `2022-05-06T00:40:21+02:00`, `contained text `+common.StaticContent(101)),
		common.NewFile(`external\unique.txt`, `2022-05-06T00:40:21+02:00`, `unique text2 `+common.StaticContent(101)),
		common.NewFile(`external\duplicate1.txt`, `2022-05-06T00:40:21+02:00`, `duplicate text `+common.StaticContent(101)),
		common.NewFile(`external\duplicate2.txt`, `2022-05-06T00:40:21+02:00`, `duplicate text `+common.StaticContent(101)),
		common.NewFile(`external\contained1.txt`, `2022-05-06T00:40:21+02:00`, `contained text `+common.StaticContent(101)),
		common.NewFile(`external\contained2.txt`, `2022-05-06T00:40:21+02:00`, `contained text `+common.StaticContent(101)),
	})
	baseDir := filepath.Join(dir, `base`)
	externalDir := filepath.Join(dir, `external`)

	common.CreateIntegrityFile(t, baseDir, []common.FileHash{
		common.NewFileHash(`c23f4692194684cb5e33e17ae1ca9290e2eeb700a148fc115b8093eb57a13cb4`, `2022-05-06T00:40:21+02:00`, `2022-05-06T00:40:21+02:00`, `103437`, `a.txt`),
		common.NewFileHash(`b9fbc96548aca1dccc257d3a8db6cda75d2a6de606e34caa77b1a3911815b625`, `2022-05-06T00:40:21+02:00`, `2022-05-06T00:40:21+02:00`, `103439`, `b.txt`),
	})

	// Execute without deletion
	fileintegrity.CheckContained(baseDir, externalDir, false, fileintegrity.EnabledOptions())

	common.AssertFilesExist(t, dir, files)
	common.AssertContainedLogFile(t, baseDir, []common.LogBlock{
		common.NewContainedLogBlock(`b9fbc96548aca1dccc257d3a8db6cda75d2a6de606e34caa77b1a3911815b625`, []string{
			`contained1.txt`,
			`contained2.txt`,
		}),
		common.NewDuplicateLogBlock(`007b935d6617881943bf90dc68055cfdf9aeb4722bd8538b0fd8cb64d66a1a03`, []string{
			`duplicate1.txt`,
			`duplicate2.txt`,
		}),
	}, 2, 1, 2)

	// Execute with deletion
	fileintegrity.CheckContained(baseDir, externalDir, true, fileintegrity.EnabledOptions())
	common.AssertFilesExist(t, dir, files[:4])
}

func TestDemoFlow(t *testing.T) {
	dir, _ := common.CreateScenario("demo", common.Files{
		common.NewFile(`images\2020 Yellowstone National Park\IMG_0091.jpg`, `2020-05-06T13:40:00+00:00`, common.StaticContent(5120)),
		common.NewFile(`images\2020 Yellowstone National Park\IMG_0137.jpg`, `2022-05-06T14:15:00+00:00`, common.StaticContent(8944)),
		common.NewFile(`images\2020 Yellowstone National Park\IMG_0235.jpg`, `2022-05-06T15:20:00+00:00`, common.StaticContent(4585)),
		common.NewFile(`images\2021 Las Vegas\DSC_0203.jpg`, `2021-05-08T21:32:00+00:00`, common.StaticContent(7815)),
		common.NewFile(`images\2021 Las Vegas\DSC_0209.jpg`, `2022-05-08T21:49:00+00:00`, common.StaticContent(5587)),
		common.NewFile(`images\2021 Las Vegas\DSC_0255.jpg`, `2022-05-08T22:55:00+00:00`, common.StaticContent(6397)),
	})

	fileintegrity.Upsert(dir, fileintegrity.EnabledOptions())

	common.UpdateFile(dir, `images\2020 Yellowstone National Park\IMG_0091.jpg`, common.StaticContent(5121), `2023-01-01T20:00:00+00:00`)

	fileintegrity.Upsert(dir, fileintegrity.EnabledOptions())
	fileintegrity.Verify(dir, fileintegrity.EnabledOptions())
}

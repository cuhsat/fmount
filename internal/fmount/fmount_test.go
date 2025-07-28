package fmount

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hiforensics/fmount/internal/posix"
	"github.com/hiforensics/utils/pkg/sys"
	"github.com/hiforensics/utils/pkg/zip"
)

func TestMain(m *testing.M) {
	sys.Progress = nil

	if _, ci := os.LookupEnv("CI"); !ci {
		os.Exit(m.Run())
	}
}

func TestMount(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode")
	}

	cases := []struct {
		name, file, path string
	}{
		{
			name: "Test mount with disk image",
			file: filepath.Join("..", "testdata", "windows.zip"),
			path: "disk",
		},
	}

	for _, tt := range cases {
		tmp, err := os.MkdirTemp(os.TempDir(), tt.path)

		if err != nil {
			t.Fatal(err)
		}

		mnt, err := os.MkdirTemp(os.TempDir(), tt.path+"-mnt")

		if err != nil {
			t.Fatal(err)
		}

		err = zip.Unzip(tt.file, tmp)

		if err != nil {
			t.Fatal(err)
		}

		t.Run(tt.name, func(t *testing.T) {
			img := filepath.Join(tmp, posix.BaseFile(tt.file))

			p, err := Mount(img, mnt, "", true, []string{})

			if err != nil {
				t.Fatal(err)
			}

			if len(p) != 1 {
				t.Fatal("partition count differs")
			}

			sys := filepath.Join(mnt, "p1")

			if p[0] != sys {
				t.Fatal("mount point does not exist")
			}

			dir, _ := os.ReadDir(sys)

			if len(dir) == 0 {
				t.Fatal("mount point is empty")
			}

			err = Unmount(img)

			if err != nil {
				t.Fatal(err)
			}

			if _, err = os.Stat(mnt); !os.IsNotExist(err) {
				t.Fatal("mount point not removed")
			}
		})
	}
}

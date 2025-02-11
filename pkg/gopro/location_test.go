package gopro

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/posener/gitfs"
	"github.com/posener/gitfs/fsutil"
	"github.com/stretchr/testify/require"
)

func TestParseGPMF(t *testing.T) {
	ctx := context.Background()

	fs, err := gitfs.New(ctx, "github.com/gopro/gpmf-parser/samples",
		gitfs.OptGlob("hero[5..6]*.mp4"))
	require.NoError(t, err)

	walk := fsutil.Walk(fs, "")
	for walk.Step() {
		if walk.Path() == "" {
			walk.Step()
		}
		fmt.Printf("\tTesting [%s]", walk.Path())

		remoteFile, err := fs.Open(walk.Path())
		require.NoError(t, err)
		localFile, err := ioutil.TempFile(".", walk.Path())
		require.NoError(t, err)
		defer os.Remove(localFile.Name())

		stat, err := remoteFile.Stat()
		require.NoError(t, err)

		buf := make([]byte, stat.Size()) // roughly 290 bytes per SRT entry

		_, err = remoteFile.Read(buf)
		require.NoError(t, err)
		require.NotEmpty(t, buf)

		_, err = localFile.Write(buf)
		require.NoError(t, err)
		err = localFile.Close()
		require.NoError(t, err)

		returned, err := fromMP4(localFile.Name())
		require.NoError(t, err)
		require.NotZero(t, returned.Latitude)
		require.NotZero(t, returned.Longitude)
		require.NotEqual(t, returned.Latitude, returned.Longitude)
		fmt.Printf("\n\treturned: %f %f\n", returned.Latitude, returned.Longitude)
	}
}

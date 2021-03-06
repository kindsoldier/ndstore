package fsfile

import (
    "bytes"
    "math/rand"
    "testing"
    "ndstore/fstore/fssrv/fsreg"
    "github.com/stretchr/testify/require"
    "github.com/stretchr/testify/assert"

)

func Test_File_WriteRead(t *testing.T) {
    var err error

    baseDir := t.TempDir()

    dbPath := "postgres://test@localhost/test"
    reg := fsreg.NewReg()

    err = reg.OpenDB(dbPath)
    require.NoError(t, err)

    err = reg.MigrateDB()
    require.NoError(t, err)

    // Set file parameters
    var batchSize   int64   = 5
    var blockSize   int64   = 1024


    // Create new file
    fileId, file0, err := NewFile(reg, baseDir, batchSize, blockSize)
    require.NoError(t, err)
    require.NotEqual(t, file0, nil)

    err = file0.Erase()
    require.NoError(t, err)

    err = file0.Close()
    require.NoError(t, err)

    // Erase all files
    fileDescrs, err := reg.ListAllFileDescrs()
    require.NoError(t, err)
    for _, fileDescr := range fileDescrs {
        file, err := OpenFile(reg, baseDir, fileDescr.FileId)
        if err != nil {
            continue
        }
        if file != nil {
            file.Delete()
            file.Close()
        }
    }

    cleanAllUnised(t, reg, baseDir)


    // Prepare data
    var dataSize int64 = batchSize * blockSize * 10 + 1
    data := make([]byte, dataSize)
    rand.Read(data)

    // Create new file
    fileId, file0, err = NewFile(reg, baseDir, batchSize, blockSize)
    require.NoError(t, err)
    require.NotEqual(t, file0, nil)
    // Write to file
    need := blockSize + 1
    count := 99
    for i := 0; i < count; i++ {
        reader0 := bytes.NewReader(data)
        written0, err := file0.Write(reader0, need)
        require.NoError(t, err)
        require.Equal(t, need, written0)
    }

    // Delete file
    //err = file0.Delete()
    //require.NoError(t, err)

    // Close file
    err = file0.Close()
    require.NoError(t, err)

    cleanAllUnised(t, reg, baseDir)

    // Reopen file
    file1, err := OpenFile(reg, baseDir, fileId)
    require.NoError(t, err)

    // Read file
    writer1 := bytes.NewBuffer(make([]byte, 0))
    read1, err := file1.Read(writer1)
    require.NoError(t, err)

    require.Equal(t, need * int64(count), int64(len(writer1.Bytes())))
    require.Equal(t, need * int64(count), read1)
    require.Equal(t, data[0:need], writer1.Bytes()[0:need])
    // Check data
    for i := 0; i < count; i++ {
        offset := int64(i) * need
        require.Equal(t, data[0:need], writer1.Bytes()[0+offset:need+offset])
    }
    // Delete file
    err = file1.Delete()
    require.NoError(t, err)

    // Delete file again
    err = file1.Delete()
    require.NoError(t, err)

    // Close file
    err = file1.Close()
    require.NoError(t, err)

    // Open non exist file
    _, err = OpenFile(reg, baseDir, fileId)
    require.Error(t, err)

    cleanAllUnised(t, reg, baseDir)
}


func cleanAllUnised(t *testing.T, reg *fsreg.Reg, baseDir string) {

    // Clean all unised files
    for {
        exists, descr, err := reg.GetAnyUnusedFileDescr()
        require.NoError(t, err)
        if !exists {
            break
        }
        file, err := OpenSpecUnusedFile(reg, baseDir, descr.FileId, descr.FileVer)
        require.NoError(t, err)
        err = file.Erase()
        require.NoError(t, err)
        err = file.Close()
        require.NoError(t, err)
    }
    // Clean all unised batchs
    for {
        exists, descr, err := reg.GetAnyUnusedBatchDescr()
        require.NoError(t, err)
        if !exists {
            break
        }
        block, err := OpenSpecUnusedBatch(reg, baseDir, descr.FileId, descr.BatchId, descr.BatchVer)
        require.NoError(t, err)
        err = block.Erase()
        require.NoError(t, err)
        err = block.Close()
        require.NoError(t, err)
    }
    // Clean all unised blocks
    for {
        exists, descr, err := reg.GetAnyUnusedBlockDescr()
        require.NoError(t, err)
        if !exists {
            break
        }
        block, err := OpenSpecUnusedBlock(reg, baseDir, descr.FileId, descr.BatchId, descr.BlockId,
                                                        descr.BlockType, descr.BlockVer)
        require.NoError(t, err)
        err = block.Erase()
        require.NoError(t, err)
        err = block.Close()
        require.NoError(t, err)
    }
}

func Benchmark_File_Write(b *testing.B) {
    // Prepare env
    var err error
    baseDir := b.TempDir()
    dbPath := "postgres://test@localhost/test"
    reg := fsreg.NewReg()

    err = reg.OpenDB(dbPath)
    require.NoError(b, err)

    err = reg.MigrateDB()
    require.NoError(b, err)

    // Erase all files
    fileDescrs, err := reg.ListAllFileDescrs()
    require.NoError(b, err)
    for _, fileDescr := range fileDescrs {
        file, err := OpenFile(reg, baseDir, fileDescr.FileId)
        if err != nil {
            continue
        }
        if file != nil {
            file.Delete()
            file.Close()
        }
    }
    // Set file parameters
    var batchSize   int64   = 5
    var blockSize   int64   = 16 * 1024 * 1024
    // Prepare data
    var dataSize int64 = (batchSize * blockSize * 10) / 16 + 1
    data := make([]byte, dataSize)
    rand.Read(data)

    need := dataSize

    b.ResetTimer()
    pBench := func(pb *testing.PB) {
        for pb.Next() {
            var err error
            // Create new file
            _, file, err := NewFile(reg, baseDir, batchSize, blockSize)
            assert.NoError(b, err)
            assert.NotEqual(b, file, nil)
            // Write to file
            reader := bytes.NewReader(data)
            written, err := file.Write(reader, need)
            assert.NoError(b, err)
            assert.Equal(b, need, written)
            // Close file
            err = file.Close()
            assert.NoError(b, err)
        }
    }
    b.SetParallelism(10)
    b.RunParallel(pBench)
}

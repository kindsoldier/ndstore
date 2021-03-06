
/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package bsapi

const LoadBlockMethod string = "loadBlock"

type LoadBlockParams struct {
    FileId      int64           `json:"fileId"`
    BatchId     int64           `json:"batchId"`
    BlockId     int64           `json:"blockId"`
    BlockType   string          `json:"blockType"`
    BlockVer    int64           `json:"blockVer"`
}

type LoadBlockResult struct {
}

func NewLoadBlockResult() *LoadBlockResult {
    return &LoadBlockResult{}
}
func NewLoadBlockParams() *LoadBlockParams {
    return &LoadBlockParams{}
}

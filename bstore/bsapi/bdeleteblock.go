
/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package bsapi

const DeleteBlockMethod string = "deleteBlock"

type DeleteBlockParams struct {
    FileId      int64           `json:"fileId"`
    BatchId     int64           `json:"batchId"`
    BlockId     int64           `json:"blockId"`
    BlockType   string          `json:"blockType"`
    BlockVer    int64           `json:"blockVer"`
}

type DeleteBlockResult struct {
}

func NewDeleteBlockResult() *DeleteBlockResult {
    return &DeleteBlockResult{}
}
func NewDeleteBlockParams() *DeleteBlockParams {
    return &DeleteBlockParams{}
}

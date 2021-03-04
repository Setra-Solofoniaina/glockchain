/*
*
* copyright - Setra Solofoniaina - 2021
 */

package main

import (
	"github.com/Setra-Solofoniaina/glockchain/tests"
	"github.com/Setra-Solofoniaina/glockchain/utils"
)

func main() {
	utils.CreateDBFolder()
	tests.MerkleTreeTest()
}

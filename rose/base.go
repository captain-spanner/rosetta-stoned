package rose

import (
	"fmt"
)

const (
	Copyright string	= "Copyright Bruce Ellis 2014"
	Version string		= "Pusbox 0.0"
	Config string		= "stone.conf"
)

const (
	hError		= hashc(iota)
	hHashed
	hIndexed
	hLiteral
)

const (
	pNone		= partc(iota)
	pCl
	pDj
	pDn
	pDr
	pDv
	pIj
	pIn
	pIr
	pIv
	pVs
	pVx
	pXj
	pXn
	pXr
	pXv
	pMax
)

var (
	root	string
	base	*corpus

	debug bool		= false
	message bool		= true
	verbose bool		= true
	xeq bool		= false
	interactive		= false
	prompt			= ">> "

	hashes map[hashc]string	= map[hashc]string {
		hError:		"error",
		hHashed:	"hashed",
		hIndexed:	"indexed",
		hLiteral:	"literal",
	}

	partm map[string]partc	= map[string]partc {
		"adj.exc":	pXj,
		"adv.exc":	pXr,
		"cntlist.rev":	pCl,
		"data.adj":	pDj,
		"data.adv":	pDr,
		"data.noun":	pDn,
		"data.verb":	pDv,
		"index.adj":	pIj,
		"index.adv":	pIr,
		"index.noun":	pIn,
		"index.verb":	pIv,
		"noun.exc":	pXn,
		"sentidx.vrb":	pVx,
		"sents.vrb":	pVs,
		"verb.exc":	pXv,
	}

	parts map[partc]string	= map[partc]string {
		pCl:	"Cl",
		pDj:	"Jd",
		pDn:	"Nd",
		pDr:	"Rd",
		pDv:	"Vd",
		pIj:	"Ji",
		pIn:	"Ni",
		pIr:	"Ri",
		pIv:	"Vi",
		pVs:	"Vs",
		pVx:	"Vx",
		pXj:	"Jx",
		pXn:	"Nx",
		pXr:	"Rx",
		pXv:	"Vx",
	}
)

type hashc int
type partc int

func init() {
	if debug {
		fmt.Println("Init...")
	}
	init_syms()
	configure()
}

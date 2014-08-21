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
	WIT		= partd(iota)
	Adjective
	Adverb
	Noun
	Verb
)

const (
	pNone		= partc(iota)
	pCl
	pCr
	pDj
	pDn
	pDr
	pDv
	pFs
	pFv
	pIj
	pIn
	pIr
	pIs
	pIv
	pVn
	pVs
	pVx
	pXj
	pXn
	pXr
	pXv
	pMax
)

const (
	dNone		= psd(iota)

	first_Adverb

	rAnt
	rDer
	rDsr
	rDst
	rDsu

	last_Adverb

	first_Adjective

	jAnt
	jAtt
	jDsr
	jDst
	jDsu
	JPar
	jPer
	jSee
	jSim

	last_Adjective

	first_Noun

	nAnt
	nAtt
	nDer
	nDsr
	nDst
	nDsu
	nHie
	nHom
	nHop
	nHos
	nHpe
	nMdr
	nMdt
	nMdu
	nMem
	nMep
	nMes
	nNio
	nNpo

	last_Noun

	first_Verb

	vAnt
	vCau
	vDer
	vDom
	vDsr
	vDst
	vDsu
	vEnt
	vGRP
	vHpe
	vNpo
	vSee

	last_Verb

	dMax
)

/*

Pointer Symbol Data


Adjective

jAnt	!	Antonym
jSim	&	Similar to
JPar	<	Participle of noun
jPer	\	Pertainym - noun
jAtt	=	Attribute
jSee	^	See also
jDst	;c	Domain of synset - topic
jDsr	;r	Domain of synset - region
jDsu	;u	Domain of synset - usage

Adverb

rAnt	!	Antonym
rDer	\	Derived from adjective
rDst	;c	Domain of synset - topic
rDsr	;r	Domain of synset - region
rDsu	;u	Domain of synset - usage
Noun

Noun

nAnt	!	Antonym
nHpe	@	Hypernm
nHie	@i	Instance Hypernm
nNpo	~	Hyponym
nNio	~i	Instance Hyponym
nHom	#m	Member holonym
nHos	#s	Substance holonym
nHop	#p	Part holonym
nMem	%m	Member meronym
nMes	%s	Substance meronym
nMep	%p	Part meronym
nAtt	=	Attribute
nDer	+	Derivationally related form
nDst	;c	Domain of synset - topic
nMdt	-c	Member of this domain - topic
nDsr	;r	Domain of synset - region
nMdr	-r	Member of this domain - region
nDsu	;u	Domain of synset - usage
nMdu	-u	Member of this domain - usage

Verb

vAnt	!	Antonym
vHpe	@	Hypernm
vNpo	~	Hyponym
vEnt	*	Entailment
vCau	>	Cause
vSee	^	See also
vGRP	$	Verb Group
vDer	+	Derivationally related form
vDst	;c	Domain of synset - topic
vDsr	;r	Domain of synset - region
vDsu	;u	Domain of synset - usage
// undocumented. assumed domain*
Vdom	;	Domain of synset
*/

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

	posm map[byte]partd	= map[byte]partd {
		'a':	Adjective,
		'r':	Adverb,
		'n':	Noun,
		'v':	Verb,
		's':	Adjective,	// what is this?
		'j':	Adjective,	// maybe not needed
	}

	poss map[partd]string	= map[partd]string {
		WIT:		"What?",
		Adjective:	"Adjective",
		Adverb:		"Adverb",
		Noun:		"Noun",
		Verb:		"Verb",
	}

	partm map[string]partc	= map[string]partc {
		// Princeton wordnet
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

		// extras from geo
		"cntlist":	pCr,
		"frames.vrb":	pFv,
		"index.sense":	pIs,
		"verb.Framestext": pFs,
	}

	parts map[partc]string	= map[partc]string {
		pCl:	"Cl",
		pCr:	"Cr",
		pDj:	"Jd",
		pDn:	"Nd",
		pDr:	"Rd",
		pDv:	"Vd",
		pFs:	"Fs",
		pFv:	"Fv",
		pIj:	"Ji",
		pIn:	"Ni",
		pIr:	"Ri",
		pIs:	"Is",
		pIv:	"Vi",
		pVn:	"Vn",
		pVs:	"Vs",
		pXj:	"Jx",
		pXn:	"Nx",
		pXr:	"Rx",
		pXv:	"Vx",
	}

	partt map[string]partc	= map[string]partc {
		"Cl":	pCl,
		"Cr":	pCr,
		"Fs":	pFs,
		"Fv":	pFv,
		"Is":	pIs,
		"Jd":	pDj,
		"Ji":	pIj,
		"Jx":	pXj,
		"Nd":	pDn,
		"Ni":	pIn,
		"Nx":	pXn,
		"Rd":	pDr,
		"Ri":	pIr,
		"Rx":	pXr,
		"Vd":	pDv,
		"Vi":	pIv,
		"Vn":	pVn,
		"Vs":	pVs,
		"Vx":	pXv,
	}

	psds	[dMax]string = [dMax]string {
		dNone:	"GOK",

		rAnt:	"rAnt",
		rDer:	"rDer",
		rDsr:	"rDsr",
		rDst:	"rDst",
		rDsu:	"rDsu",

		jAnt:	"jAnt",
		jAtt:	"jAtt",
		jDsr:	"jDsr",
		jDst:	"jDst",
		jDsu:	"jDsu",
		JPar:	"JPar",
		jPer:	"jPer",
		jSee:	"jSee",
		jSim:	"jSim",

		nAnt:	"nAnt",
		nAtt:	"nAtt",
		nDer:	"nDer",
		nDsr:	"nDsr",
		nDst:	"nDst",
		nDsu:	"nDsu",
		nHie:	"nHie",
		nHom:	"nHom",
		nHop:	"nHop",
		nHos:	"nHos",
		nHpe:	"nHpe",
		nMdr:	"nMdr",
		nMdt:	"nMdt",
		nMdu:	"nMdu",
		nMem:	"nMem",
		nMep:	"nMep",
		nMes:	"nMes",
		nNio:	"nNio",
		nNpo:	"nNpo",

		vAnt:	"vAnt",
		vCau:	"vCau",
		vDer:	"vDer",
		vDom:	"vDom",
		vDsr:	"vDsr",
		vDst:	"vDst",
		vDsu:	"vDsu",
		vEnt:	"vEnt",
		vGRP:	"vGRP",
		vHpe:	"vHpe",
		vNpo:	"vNpo",
		vSee:	"vSee",
	}

	psdd	[dMax]string = [dMax]string {
		dNone:	"Unknown",

		rAnt:	"Antonym",
		rDer:	"Derived from adjective",
		rDst:	"Domain of synset - topic",
		rDsr:	"Domain of synset - region",
		rDsu:	"Domain of synset - usage",

		jAnt:	"Antonym",
		jSim:	"Similar to",
		JPar:	"Participle of noun",
		jPer:	"Pertainym - noun",
		jAtt:	"Attribute",
		jSee:	"See also",
		jDst:	"Domain of synset - topic",
		jDsr:	"Domain of synset - region",
		jDsu:	"Domain of synset - usage",

		nAnt:	"Antonym",
		nHpe:	"Hypernm",
		nHie:	"Instance Hypernm",
		nNpo:	"Hyponym",
		nNio:	"Instance Hyponym",
		nHom:	"Member holonym",
		nHos:	"Substance holonym",
		nHop:	"Part holonym",
		nMem:	"Member meronym",
		nMes:	"Substance meronym",
		nMep:	"Part meronym",
		nAtt:	"Attribute",
		nDer:	"Derivationally related form",
		nDst:	"Domain of synset - topic",
		nMdt:	"Member of this domain - topic",
		nDsr:	"Domain of synset - region",
		nMdr:	"Member of this domain - region",
		nDsu:	"Domain of synset - usage",
		nMdu:	"Member of this domain - usage",

		vAnt:	"Antonym",
		vHpe:	"Hypernm",
		vNpo:	"Hyponym",
		vEnt:	"Entailment",
		vCau:	"Cause",
		vSee:	"See also",
		vGRP:	"Verb Group",
		vDer:	"Derivationally related form",
		vDom:	"Domain of synset",
		vDst:	"Domain of synset - topic",
		vDsr:	"Domain of synset - region",
		vDsu:	"Domain of synset - usage",
	}

	adjpsdm	map[string]psd = map[string]psd {
		"!":	jAnt,
		"&":	jSim,
		"<":	JPar,
		"\\":	jPer,
		"=":	jAtt,
		"^":	jSee,
		";c":	jDst,
		";r":	jDsr,
		";u":	jDsu,
	}

	advpsdm	map[string]psd = map[string]psd {
		"!":	rAnt,
		"\\":	rDer,
		";c":	rDst,
		";r":	rDsr,
		";u":	rDsu,
	}

	nounpsdm	map[string]psd = map[string]psd {
		"!":	nAnt,
		"@":	nHpe,
		"@i":	nHie,
		"~":	nNpo,
		"~i":	nNio,
		"#m":	nHom,
		"#s":	nHos,
		"#p":	nHop,
		"%m":	nMem,
		"%s":	nMes,
		"%p":	nMep,
		"=":	nAtt,
		"+":	nDer,
		";c":	nDst,
		"-c":	nMdt,
		";r":	nDsr,
		"-r":	nMdr,
		";u":	nDsu,
		"-u":	nMdu,
	}

	verbpsdm	map[string]psd = map[string]psd {
		"!":	vAnt,
		"@":	vHpe,
		"~":	vNpo,
		"*":	vEnt,
		">":	vCau,
		"^":	vSee,
		"$":	vGRP,
		"+":	vDer,
		";":	vDom,
		";c":	vDst,
		";r":	vDsr,
		";u":	vDsu,
	}
	
	posdx	[]string = []string {
		Adjective:	"Jd",
		Adverb:		"Rd",
		Noun:		"Nd",
		Verb:		"Vd",
	}
	
	posmv	[]map[string]psd = []map[string]psd {
		Adjective:	adjpsdm,
		Adverb:		advpsdm,
		Noun:		nounpsdm,
		Verb:		verbpsdm,
	}

	psdmv	[]map[string]psd = []map[string]psd {
		pIj:	adjpsdm,
		pIr:	advpsdm,
		pIn:	nounpsdm,
		pIv:	verbpsdm,
		pDj:	adjpsdm,
		pDr:	advpsdm,
		pDn:	nounpsdm,
		pDv:	verbpsdm,
	}
)

type hashc int
type partc int
type partd byte
type psd byte

func init() {
	if debug {
		fmt.Println("Init...")
	}
	go petalsrv()
	init_syms()
	init_cmds()
	configure()
}

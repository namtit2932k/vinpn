package formats

import (
	"bytes"

	"github.com/sickyturtlez/vinpn/internal/rules"
) // VinPNHeader is the first line of a VinPN rules list.
const VinPNHeader = "# vinpn-rules v1"

// LegacyFormat and LegacyHeader are the pre-rename format id and list
// header. They stay readable because signed lists in the wild still ship
// the old header (editing them would break their .sig) and imported
// backups from an old install carry the old id.
const LegacyFormat Format = "ghostline"

// LegacyHeader is the pre-rename list header.
const LegacyHeader = "# ghostline-rules v1"

// HasVinPNHeader reports whether data is a VinPN-format rules list by its
// first non-blank line (current or legacy header).
func HasVinPNHeader(data []byte) bool {
	data = bytes.TrimLeft(bytes.TrimPrefix(data, []byte{0xEF, 0xBB, 0xBF}), " \t\r\n")
	line, _, _ := bytes.Cut(data, []byte("\n"))
	h := string(bytes.TrimSpace(line))
	return h == VinPNHeader || h == LegacyHeader
}

func isVinPN(data []byte) bool { return HasVinPNHeader(data) }

// parseVinPN reads one rule-text line with its own actions. Lines the
// rule parser rejects (sni= on a keyword, unknown upstreams…) are skipped.
func parseVinPN(b *builder, line string, n int) {
	if isComment(line) {
		return
	}
	rs, errs := rules.ParseText(line, nil)
	if len(errs) > 0 || len(rs) != 1 {
		b.skip(line)
		return
	}
	p, err := rules.ParsePattern(rs[0].Pattern)
	if err != nil {
		b.skip(line)
		return
	}
	a := rs[0].Action
	b.add(p, n, line, false, nil)
	b.r.Entries[len(b.r.Entries)-1].Action = &a
}

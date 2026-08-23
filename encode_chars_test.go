package iabtcfv2

import (
	"testing"
)

// BACK-853 : WriteChars écrit n/bitsChar caractères en indexant v[i] sans borne.
// Une valeur plus courte que le champ — ConsentLanguage ou PublisherCC vide ou
// d'un seul caractère — provoque donc un dépassement d'indice. PublisherCC vide
// étant la valeur par défaut de CoreString, le cas se produit sur le chemin le
// plus simple : construire un CoreString sans renseigner ce champ.
//
// Le champ fait 12 bits en toutes circonstances. La valeur est complétée à
// droite par des zéros lorsqu'elle est plus courte, et tronquée lorsqu'elle est
// plus longue — c'est le comportement de la plateforme de gestion du
// consentement de la maison, qui produit les chaînes réellement servies.
func TestWriteCharsAcceptsValuesShorterThanTheField(t *testing.T) {
	cases := []struct {
		name  string
		value string
		want  string
	}{
		{"deux caractères, cas nominal", "FR", "FR"},
		{"valeur vide, défaut de la structure", "", "AA"},
		{"un seul caractère", "F", "FA"},
		{"plus longue que le champ, tronquée", "FRA", "FR"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			e := NewTCEncoder(make([]byte, bitsConsentLanguage/8+1))
			e.WriteChars(c.value, bitsConsentLanguage)

			e.Position = 0
			if got := e.ReadChars(bitsConsentLanguage); got != c.want {
				t.Errorf("WriteChars(%q) puis ReadChars = %q, attendu %q", c.value, got, c.want)
			}
			if e.Position != bitsConsentLanguage {
				t.Errorf("WriteChars(%q) : %d bits relus, le champ en fait %d", c.value, e.Position, bitsConsentLanguage)
			}
		})
	}
}

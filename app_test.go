package main

import (
	"testing"

	"language-refresher/internal/settings"
)

func TestNormalizeLocale(t *testing.T) {
	cases := map[string]string{
		"pt-BR": "pt-BR",
		"pt-PT": "pt-BR",
		"PT_br": "pt-BR",
		"en-US": "en",
		"es":    "en",
		"":      "en",
	}
	for input, expected := range cases {
		if actual := normalizeLocale(input); actual != expected {
			t.Errorf("normalizeLocale(%q) = %q; esperado %q", input, actual, expected)
		}
	}
}

func TestResolveLocalePrefersSavedChoice(t *testing.T) {
	if locale := resolveLocale("pt-BR", settings.State{Language: "en"}, true); locale != "en" {
		t.Fatalf("escolha salva deveria prevalecer, recebeu %q", locale)
	}
	if locale := resolveLocale("en-US", settings.State{}, false); locale != "en" {
		t.Fatalf("localidade do sistema deveria ser detectada, recebeu %q", locale)
	}
}

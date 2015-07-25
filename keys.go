package main

import key "golang.org/x/mobile/event/key"

type KeyBinding uint8

const (
	KeyUnknown KeyBinding = iota
	KeyForward
	KeyReverse
	KeyStrafeLeft
	KeyStrafeRight
	KeyLineLeft
	KeyLineRight
	KeyCameraFollow
	KeyPause
)

func DefaultBindings() *Bindings {
	binding := map[key.Code]KeyBinding{
		key.CodeW:          KeyForward,
		key.CodeS:          KeyReverse,
		key.CodeA:          KeyStrafeLeft,
		key.CodeD:          KeyStrafeRight,
		key.CodeRightArrow: KeyLineRight,
		key.CodeLeftArrow:  KeyLineLeft,
		key.CodeF:          KeyCameraFollow,
		key.CodeSpacebar:   KeyPause,
	}
	return &Bindings{
		bindings: binding,
		on:       map[KeyBinding]func(KeyBinding){},
		pressed:  map[KeyBinding]bool{},
	}
}

type Bindings struct {
	bindings map[key.Code]KeyBinding
	on       map[KeyBinding]func(KeyBinding)
	pressed  map[KeyBinding]bool
}

func (b *Bindings) Lookup(code key.Code) KeyBinding {
	k, ok := b.bindings[code]
	if !ok {
		return KeyUnknown
	}
	return k
}

func (b *Bindings) On(k KeyBinding, fn func(KeyBinding)) {
	b.on[k] = fn
}

func (b *Bindings) Press(code key.Code) {
	key := b.Lookup(code)
	b.pressed[key] = true

	fn, ok := b.on[key]
	if ok {
		fn(key)
	}
}

func (b *Bindings) Release(code key.Code) {
	key := b.Lookup(code)
	b.pressed[key] = false
}

func (b *Bindings) Pressed(k KeyBinding) bool {
	p, ok := b.pressed[k]
	if ok {
		return p
	}
	return false
}

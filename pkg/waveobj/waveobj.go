// Copyright 2024, Command Line Inc.
// SPDX-License-Identifier: Apache-2.0

// Package waveobj defines the core object model for WaveTerm.
// All persistent objects in WaveTerm implement the WaveObj interface.
package waveobj

import (
	"fmt"
	"reflect"
)

// OType is a string identifier for a WaveObj type (e.g. "tab", "block", "workspace")
type OType = string

// OID is a unique object identifier (UUID string)
type OID = string

// MetaMapType is a generic map for storing metadata on WaveObjs
type MetaMapType = map[string]any

// WaveObj is the base interface for all persistent WaveTerm objects.
// Every object must expose its type, unique ID, and version for OCC (optimistic concurrency control).
type WaveObj interface {
	GetOType() OType
	GetOID() OID
	GetVersion() int
	SetVersion(v int)
}

// WaveObjBase provides a default implementation of WaveObj fields
// that can be embedded in concrete types.
type WaveObjBase struct {
	OType   OType  `json:"otype"`
	OID     OID    `json:"oid"`
	Version int    `json:"version"`
}

func (b *WaveObjBase) GetOType() OType   { return b.OType }
func (b *WaveObjBase) GetOID() OID       { return b.OID }
func (b *WaveObjBase) GetVersion() int   { return b.Version }
func (b *WaveObjBase) SetVersion(v int)  { b.Version = v }

// typeRegistry maps OType strings to their reflect.Type for dynamic instantiation.
var typeRegistry = map[OType]reflect.Type{}

// RegisterType registers a WaveObj implementation type under the given OType name.
// This must be called during init() for each concrete WaveObj type.
func RegisterType(otype OType, rtype reflect.Type) {
	if _, exists := typeRegistry[otype]; exists {
		panic(fmt.Sprintf("waveobj: type %q already registered", otype))
	}
	if rtype.Kind() == reflect.Ptr {
		rtype = rtype.Elem()
	}
	typeRegistry[otype] = rtype
}

// MakeWaveObj creates a new zero-value WaveObj instance for the given OType.
// Returns an error if the type has not been registered.
func MakeWaveObj(otype OType) (WaveObj, error) {
	rtype, ok := typeRegistry[otype]
	if !ok {
		return nil, fmt.Errorf("waveobj: unknown otype %q", otype)
	}
	obj, ok := reflect.New(rtype).Interface().(WaveObj)
	if !ok {
		return nil, fmt.Errorf("waveobj: registered type %q does not implement WaveObj", otype)
	}
	return obj, nil
}

// GetOTypeOf returns the OType for a registered reflect.Type, or empty string if not found.
func GetOTypeOf(rtype reflect.Type) OType {
	if rtype.Kind() == reflect.Ptr {
		rtype = rtype.Elem()
	}
	for otype, rt := range typeRegistry {
		if rt == rtype {
			return otype
		}
	}
	return ""
}

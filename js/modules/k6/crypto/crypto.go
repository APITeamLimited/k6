/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2017 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package crypto

import (
	"context"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"hash"

	"golang.org/x/crypto/md4"
	"golang.org/x/crypto/ripemd160"

	"github.com/dop251/goja"

	"go.k6.io/k6/js/common"
)

type Crypto struct***REMOVED******REMOVED***

type Hasher struct ***REMOVED***
	ctx context.Context

	hash hash.Hash
***REMOVED***

func New() *Crypto ***REMOVED***
	return &Crypto***REMOVED******REMOVED***
***REMOVED***

// RandomBytes returns random data of the given size.
func (*Crypto) RandomBytes(ctx context.Context, size int) *goja.ArrayBuffer ***REMOVED***
	rt := common.GetRuntime(ctx)
	if size < 1 ***REMOVED***
		common.Throw(rt, errors.New("invalid size"))
	***REMOVED***
	bytes := make([]byte, size)
	_, err := rand.Read(bytes)
	if err != nil ***REMOVED***
		common.Throw(rt, err)
	***REMOVED***
	ab := rt.NewArrayBuffer(bytes)
	return &ab
***REMOVED***

// Md4 returns the MD4 hash of input in the given encoding.
func (c *Crypto) Md4(ctx context.Context, input interface***REMOVED******REMOVED***, outputEncoding string) interface***REMOVED******REMOVED*** ***REMOVED***
	hasher := c.CreateHash(ctx, "md4")
	hasher.Update(input)
	return hasher.Digest(outputEncoding)
***REMOVED***

// Md5 returns the MD5 hash of input in the given encoding.
func (c *Crypto) Md5(ctx context.Context, input interface***REMOVED******REMOVED***, outputEncoding string) interface***REMOVED******REMOVED*** ***REMOVED***
	hasher := c.CreateHash(ctx, "md5")
	hasher.Update(input)
	return hasher.Digest(outputEncoding)
***REMOVED***

// Sha1 returns the SHA1 hash of input in the given encoding.
func (c *Crypto) Sha1(ctx context.Context, input interface***REMOVED******REMOVED***, outputEncoding string) interface***REMOVED******REMOVED*** ***REMOVED***
	hasher := c.CreateHash(ctx, "sha1")
	hasher.Update(input)
	return hasher.Digest(outputEncoding)
***REMOVED***

// Sha256 returns the SHA256 hash of input in the given encoding.
func (c *Crypto) Sha256(ctx context.Context, input interface***REMOVED******REMOVED***, outputEncoding string) interface***REMOVED******REMOVED*** ***REMOVED***
	hasher := c.CreateHash(ctx, "sha256")
	hasher.Update(input)
	return hasher.Digest(outputEncoding)
***REMOVED***

// Sha384 returns the SHA384 hash of input in the given encoding.
func (c *Crypto) Sha384(ctx context.Context, input interface***REMOVED******REMOVED***, outputEncoding string) interface***REMOVED******REMOVED*** ***REMOVED***
	hasher := c.CreateHash(ctx, "sha384")
	hasher.Update(input)
	return hasher.Digest(outputEncoding)
***REMOVED***

// Sha512 returns the SHA512 hash of input in the given encoding.
func (c *Crypto) Sha512(ctx context.Context, input interface***REMOVED******REMOVED***, outputEncoding string) interface***REMOVED******REMOVED*** ***REMOVED***
	hasher := c.CreateHash(ctx, "sha512")
	hasher.Update(input)
	return hasher.Digest(outputEncoding)
***REMOVED***

// Sha512_224 returns the SHA512/224 hash of input in the given encoding.
func (c *Crypto) Sha512_224(ctx context.Context, input interface***REMOVED******REMOVED***, outputEncoding string) interface***REMOVED******REMOVED*** ***REMOVED***
	hasher := c.CreateHash(ctx, "sha512_224")
	hasher.Update(input)
	return hasher.Digest(outputEncoding)
***REMOVED***

// Sha512_256 returns the SHA512/256 hash of input in the given encoding.
func (c *Crypto) Sha512_256(ctx context.Context, input interface***REMOVED******REMOVED***, outputEncoding string) interface***REMOVED******REMOVED*** ***REMOVED***
	hasher := c.CreateHash(ctx, "sha512_256")
	hasher.Update(input)
	return hasher.Digest(outputEncoding)
***REMOVED***

// Ripemd160 returns the RIPEMD160 hash of input in the given encoding.
func (c *Crypto) Ripemd160(ctx context.Context, input interface***REMOVED******REMOVED***, outputEncoding string) interface***REMOVED******REMOVED*** ***REMOVED***
	hasher := c.CreateHash(ctx, "ripemd160")
	hasher.Update(input)
	return hasher.Digest(outputEncoding)
***REMOVED***

// CreateHash returns a Hasher instance that uses the given algorithm.
func (*Crypto) CreateHash(ctx context.Context, algorithm string) *Hasher ***REMOVED***
	hasher := Hasher***REMOVED******REMOVED***
	hasher.ctx = ctx

	switch algorithm ***REMOVED***
	case "md4":
		hasher.hash = md4.New()
	case "md5":
		hasher.hash = md5.New()
	case "sha1":
		hasher.hash = sha1.New()
	case "sha256":
		hasher.hash = sha256.New()
	case "sha384":
		hasher.hash = sha512.New384()
	case "sha512_224":
		hasher.hash = sha512.New512_224()
	case "sha512_256":
		hasher.hash = sha512.New512_256()
	case "sha512":
		hasher.hash = sha512.New()
	case "ripemd160":
		hasher.hash = ripemd160.New()
	***REMOVED***

	return &hasher
***REMOVED***

// Update the hash with the input data.
func (hasher *Hasher) Update(input interface***REMOVED******REMOVED***) ***REMOVED***
	d, err := common.ToBytes(input)
	if err != nil ***REMOVED***
		common.Throw(common.GetRuntime(hasher.ctx), err)
	***REMOVED***
	_, err = hasher.hash.Write(d)
	if err != nil ***REMOVED***
		common.Throw(common.GetRuntime(hasher.ctx), err)
	***REMOVED***
***REMOVED***

// Digest returns the hash value in the given encoding.
func (hasher *Hasher) Digest(outputEncoding string) interface***REMOVED******REMOVED*** ***REMOVED***
	sum := hasher.hash.Sum(nil)
	rt := common.GetRuntime(hasher.ctx)

	switch outputEncoding ***REMOVED***
	case "base64":
		return base64.StdEncoding.EncodeToString(sum)

	case "base64url":
		return base64.URLEncoding.EncodeToString(sum)

	case "base64rawurl":
		return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(sum)

	case "hex":
		return hex.EncodeToString(sum)

	case "binary":
		ab := rt.NewArrayBuffer(sum)
		return &ab

	default:
		err := errors.New("Invalid output encoding: " + outputEncoding)
		common.Throw(rt, err)
	***REMOVED***

	return ""
***REMOVED***

// HexEncode returns a string with the hex representation of the provided byte
// array or ArrayBuffer.
func (c Crypto) HexEncode(ctx context.Context, data interface***REMOVED******REMOVED***) string ***REMOVED***
	d, err := common.ToBytes(data)
	if err != nil ***REMOVED***
		common.Throw(common.GetRuntime(ctx), err)
	***REMOVED***
	return hex.EncodeToString(d)
***REMOVED***

// CreateHMAC returns a new HMAC hash using the given algorithm and key.
func (c Crypto) CreateHMAC(ctx context.Context, algorithm string, key interface***REMOVED******REMOVED***) *Hasher ***REMOVED***
	hasher := Hasher***REMOVED******REMOVED***
	hasher.ctx = ctx
	kb, err := common.ToBytes(key)
	if err != nil ***REMOVED***
		common.Throw(common.GetRuntime(hasher.ctx), err)
	***REMOVED***

	switch algorithm ***REMOVED***
	case "md4":
		hasher.hash = hmac.New(md4.New, kb)
	case "md5":
		hasher.hash = hmac.New(md5.New, kb)
	case "sha1":
		hasher.hash = hmac.New(sha1.New, kb)
	case "sha256":
		hasher.hash = hmac.New(sha256.New, kb)
	case "sha384":
		hasher.hash = hmac.New(sha512.New384, kb)
	case "sha512_224":
		hasher.hash = hmac.New(sha512.New512_224, kb)
	case "sha512_256":
		hasher.hash = hmac.New(sha512.New512_256, kb)
	case "sha512":
		hasher.hash = hmac.New(sha512.New, kb)
	case "ripemd160":
		hasher.hash = hmac.New(ripemd160.New, kb)
	default:
		err := errors.New("Invalid algorithm: " + algorithm)
		common.Throw(common.GetRuntime(hasher.ctx), err)
	***REMOVED***

	return &hasher
***REMOVED***

// Hmac returns a new HMAC hash of input using the given algorithm and key
// in the given encoding.
func (c *Crypto) Hmac(
	ctx context.Context, algorithm string, key, input interface***REMOVED******REMOVED***, outputEncoding string,
) interface***REMOVED******REMOVED*** ***REMOVED***
	hasher := c.CreateHMAC(ctx, algorithm, key)
	hasher.Update(input)
	return hasher.Digest(outputEncoding)
***REMOVED***

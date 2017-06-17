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
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"hash"

	"golang.org/x/crypto/md4"
	"golang.org/x/crypto/ripemd160"

	"github.com/loadimpact/k6/js/common"
)

type Crypto struct***REMOVED******REMOVED***

type Hasher struct ***REMOVED***
	ctx context.Context

	hash hash.Hash
***REMOVED***

func (c *Crypto) Md4(ctx context.Context, input string, outputEncoding string) string ***REMOVED***
	hasher := c.CreateHash(ctx, "md4")
	hasher.Update(input)
	return hasher.Digest(outputEncoding)
***REMOVED***

func (c *Crypto) Md5(ctx context.Context, input string, outputEncoding string) string ***REMOVED***
	hasher := c.CreateHash(ctx, "md5")
	hasher.Update(input)
	return hasher.Digest(outputEncoding)
***REMOVED***

func (c *Crypto) Sha1(ctx context.Context, input string, outputEncoding string) string ***REMOVED***
	hasher := c.CreateHash(ctx, "sha1")
	hasher.Update(input)
	return hasher.Digest(outputEncoding)
***REMOVED***

func (c *Crypto) Sha256(ctx context.Context, input string, outputEncoding string) string ***REMOVED***
	hasher := c.CreateHash(ctx, "sha256")
	hasher.Update(input)
	return hasher.Digest(outputEncoding)
***REMOVED***

func (c *Crypto) Sha384(ctx context.Context, input string, outputEncoding string) string ***REMOVED***
	hasher := c.CreateHash(ctx, "sha384")
	hasher.Update(input)
	return hasher.Digest(outputEncoding)
***REMOVED***

func (c *Crypto) Sha512(ctx context.Context, input string, outputEncoding string) string ***REMOVED***
	hasher := c.CreateHash(ctx, "sha512")
	hasher.Update(input)
	return hasher.Digest(outputEncoding)
***REMOVED***

func (c *Crypto) Sha512_224(ctx context.Context, input string, outputEncoding string) string ***REMOVED***
	hasher := c.CreateHash(ctx, "sha512_224")
	hasher.Update(input)
	return hasher.Digest(outputEncoding)
***REMOVED***

func (c *Crypto) Sha512_256(ctx context.Context, input string, outputEncoding string) string ***REMOVED***
	hasher := c.CreateHash(ctx, "sha512_256")
	hasher.Update(input)
	return hasher.Digest(outputEncoding)
***REMOVED***

func (c *Crypto) Ripemd160(ctx context.Context, input string, outputEncoding string) string ***REMOVED***
	hasher := c.CreateHash(ctx, "ripemd160")
	hasher.Update(input)
	return hasher.Digest(outputEncoding)
***REMOVED***

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

func (hasher *Hasher) Update(input string) ***REMOVED***
	_, err := hasher.hash.Write([]byte(input))
	if err != nil ***REMOVED***
		common.Throw(common.GetRuntime(hasher.ctx), err)
	***REMOVED***
***REMOVED***

func (hasher *Hasher) Digest(outputEncoding string) string ***REMOVED***
	sum := hasher.hash.Sum(nil)

	switch outputEncoding ***REMOVED***
	case "base64":
		return base64.StdEncoding.EncodeToString(sum)

	case "hex":
		return hex.EncodeToString(sum)

	default:
		err := errors.New("Invalid output encoding: " + outputEncoding)
		common.Throw(common.GetRuntime(hasher.ctx), err)
	***REMOVED***

	return ""
***REMOVED***

func (c Crypto) CreateHmac(ctx context.Context, algorithm string, key string) *Hasher ***REMOVED***
	hasher := Hasher***REMOVED******REMOVED***
	hasher.ctx = ctx
	keyBuffer := []byte(key)

	switch algorithm ***REMOVED***
	case "md4":
		hasher.hash = hmac.New(md4.New, keyBuffer)
	case "md5":
		hasher.hash = hmac.New(md5.New, keyBuffer)
	case "sha1":
		hasher.hash = hmac.New(sha1.New, keyBuffer)
	case "sha256":
		hasher.hash = hmac.New(sha256.New, keyBuffer)
	case "sha384":
		hasher.hash = hmac.New(sha512.New384, keyBuffer)
	case "sha512_224":
		hasher.hash = hmac.New(sha512.New512_224, keyBuffer)
	case "sha512_256":
		hasher.hash = hmac.New(sha512.New512_256, keyBuffer)
	case "sha512":
		hasher.hash = hmac.New(sha512.New, keyBuffer)
	case "ripemd160":
		hasher.hash = hmac.New(ripemd160.New, keyBuffer)
	***REMOVED***

	return &hasher
***REMOVED***

func (c *Crypto) Hmac(ctx context.Context, algorithm string, key string, input string, outputEncoding string) string ***REMOVED***
	hasher := c.CreateHmac(ctx, algorithm, key)
	hasher.Update(input)
	return hasher.Digest(outputEncoding)
***REMOVED***

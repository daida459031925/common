package pbkdf2 // import "golang.org/x/crypto/pbkdf2"

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"hash"
	mathrand "math/rand"
)

const (
	saltMinLen = 8
	saltMaxLen = 32
	iter       = 1000
	keyLen     = 32
)

// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package pbkdf2 implements the key derivation function PBKDF2 as defined in RFC
2898 / PKCS #5 v2.0.

A key derivation function is useful when encrypting data based on a password
or any other not-fully-random data. It uses a pseudorandom function to derive
a secure encryption key based on the password.

While v2.0 of the standard defines only one pseudorandom function to use,
HMAC-SHA1, the drafted v2.1 specification allows use of all five FIPS Approved
Hash Functions SHA-1, SHA-224, SHA-256, SHA-384 and SHA-512 for HMAC. To
choose, you can pass the `New` functions from the different SHA packages to
pbkdf2.Key.
*/
// Key derives a key from the password, salt and iteration count, returning a
// []byte of length keylen that can be used as cryptographic key. The key is
// derived based on the method described as PBKDF2 with the HMAC variant using
// the supplied hash function.
//
// For example, to use a HMAC-SHA-1 based PBKDF2 key derivation function, you
// can get a derived key for e.g. AES-256 (which needs a 32-byte key) by
// doing:
//
// 	dk := pbkdf2.Key([]byte("some password"), salt, 4096, 32, sha1.New)
//
// Remember to get a good random salt. At least 8 bytes is recommended by the
// RFC.
//
// Using a higher iteration count will increase the cost of an exhaustive
// search but will also make derivation proportionally slower.
func key(password, salt []byte, iter, keyLen int, h func() hash.Hash) []byte {
	prf := hmac.New(h, password)
	hashLen := prf.Size()
	numBlocks := (keyLen + hashLen - 1) / hashLen

	var buf [4]byte
	dk := make([]byte, 0, numBlocks*hashLen)
	U := make([]byte, hashLen)
	for block := 1; block <= numBlocks; block++ {
		// N.B.: || means concatenation, ^ means XOR
		// for each block T_i = U_1 ^ U_2 ^ ... ^ U_iter
		// U_1 = PRF(password, salt || uint(i))
		prf.Reset()
		prf.Write(salt)
		buf[0] = byte(block >> 24)
		buf[1] = byte(block >> 16)
		buf[2] = byte(block >> 8)
		buf[3] = byte(block)
		prf.Write(buf[:4])
		dk = prf.Sum(dk)
		T := dk[len(dk)-hashLen:]
		copy(U, T)

		// U_n = PRF(password, U_(n-1))
		for n := 2; n <= iter; n++ {
			prf.Reset()
			prf.Write(U)
			U = U[:0]
			U = prf.Sum(U)
			for x := range U {
				T[x] ^= U[x]
			}
		}
	}
	return dk[:keyLen]
}

// EncryptPwd 加密密码
func EncryptPwd(pwd string) (encrypt string, salt []byte, err error) {
	// 1、生成随机长度的盐值
	salt, err = randSalt()
	if err != nil {
		return
	}

	// 2、生成加密串
	en := encryptPwdWithSalt([]byte(pwd), salt)
	en = append(en, salt...)

	// 3、合并盐值
	encrypt = base64.StdEncoding.EncodeToString(en)

	return
}

func randSalt() ([]byte, error) {
	// 生成8-32之间的随机数字
	salt := make([]byte, mathrand.Intn(saltMaxLen-saltMinLen)+saltMinLen)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

func encryptPwdWithSalt(pwd, salt []byte) (pwdEn []byte) {
	pwd = append(pwd, salt...)
	pwdEn = key(pwd, salt, iter, keyLen, sha256.New)
	return
}

// CheckEncryptPwdMatch 验证输入的密码是否与加密后字符串匹配
func CheckEncryptPwdMatch(inPwd, encrypt string) (ok bool) {
	// 1、参数校验
	if len(encrypt) == 0 {
		return
	}

	enDecode, err := base64.StdEncoding.DecodeString(encrypt)
	if err != nil {
		return
	}

	// 2、截取加密串 固定长度
	salt := enDecode[keyLen:]

	// 3、比对
	enBase64 := base64.StdEncoding.EncodeToString(enDecode[0:keyLen])
	pwdEnBase64 := base64.StdEncoding.EncodeToString(encryptPwdWithSalt([]byte(inPwd), salt))
	ok = enBase64 == pwdEnBase64

	return
}

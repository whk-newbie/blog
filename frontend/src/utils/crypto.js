import CryptoJS from 'crypto-js'

/**
 * Generate a random AES-256 key (32 bytes) as a CryptoJS WordArray
 */
export function generateAESKey() {
  return CryptoJS.lib.WordArray.random(32)
}

/**
 * AES-CBC encrypt. Returns Base64 string (IV prepended to ciphertext).
 */
export function aesEncrypt(plaintext, key) {
  const iv = CryptoJS.lib.WordArray.random(16)
  const encrypted = CryptoJS.AES.encrypt(plaintext, key, {
    iv: iv,
    mode: CryptoJS.mode.CBC,
    padding: CryptoJS.pad.Pkcs7,
  })
  const combined = iv.concat(encrypted.ciphertext)
  return CryptoJS.enc.Base64.stringify(combined)
}

/**
 * AES-CBC decrypt. Expects Base64 string with IV prepended.
 */
export function aesDecrypt(ciphertextB64, key) {
  const combined = CryptoJS.enc.Base64.parse(ciphertextB64)
  const iv = CryptoJS.lib.WordArray.create(combined.words.slice(0, 4), 16)
  const ciphertext = CryptoJS.lib.WordArray.create(
    combined.words.slice(4),
    combined.sigBytes - 16
  )
  const decrypted = CryptoJS.AES.decrypt(
    { ciphertext: ciphertext },
    key,
    {
      iv: iv,
      mode: CryptoJS.mode.CBC,
      padding: CryptoJS.pad.Pkcs7,
    }
  )
  return decrypted.toString(CryptoJS.enc.Utf8)
}

/**
 * RSA-OAEP encrypt using Web Crypto API.
 * Encrypts the AES key (as base64 string) with the server's RSA public key.
 */
export async function rsaEncrypt(publicKeyPEM, plaintext) {
  const pemContents = publicKeyPEM
    .replace('-----BEGIN PUBLIC KEY-----', '')
    .replace('-----END PUBLIC KEY-----', '')
    .replace(/\s/g, '')
  const binaryDer = Uint8Array.from(atob(pemContents), c => c.charCodeAt(0))

  const publicKey = await crypto.subtle.importKey(
    'spki',
    binaryDer.buffer,
    { name: 'RSA-OAEP', hash: 'SHA-256' },
    false,
    ['encrypt']
  )

  const encoded = new TextEncoder().encode(plaintext)
  const encrypted = await crypto.subtle.encrypt(
    { name: 'RSA-OAEP' },
    publicKey,
    encoded
  )

  return btoa(String.fromCharCode(...new Uint8Array(encrypted)))
}

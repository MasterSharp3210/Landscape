import platform
import hashlib
import secrets
import time
import json
import random
import base64
import struct
import string
import uuid
from faker import Faker

print("--------------- LandScape Encryptor ---------------")

sentiero = input("\nInserisci la PWD Iniziale: ").encode()
albero = int(input("Inserisci il numero intermedio dell' albero: "))

# Generatore stringe
def FingerPrintGenerator():
    dati_macchina = (
        str(uuid.getnode()) +
        platform.system() +
        platform.node() +
        platform.machine() +
        platform.processor()
    )

    extra = secrets.token_hex(8)

    possibili = [albero - 1, albero, albero + 1]
    scelto = random.choice(possibili)

    risultato = (dati_macchina + extra + str(scelto)).encode()

    for _ in range(100):
        risultato = hashlib.sha256(risultato).digest()

    fingerprint = hashlib.sha256(risultato).hexdigest()

    return fingerprint

# Faker delle credenziali email + password
def credentials_faker():
    locales = ['en_US', 'es_ES', 'fr_FR', 'de_DE', 'it_IT', 'ja_JP', 'ko_KR', 'zh_CN', 'ru_RU']
    
    email_providers = {
        'en_US': ['gmail.com', 'yahoo.com', 'outlook.com', 'hotmail.com', 'icloud.com'],
        'es_ES': ['gmail.com', 'yahoo.es', 'outlook.es', 'hotmail.es'],
        'fr_FR': ['gmail.com', 'yahoo.fr', 'outlook.fr', 'hotmail.fr'],
        'de_DE': ['gmail.com', 'web.de', 'gmx.de', 'outlook.de'],
        'it_IT': ['gmail.com', 'yahoo.it', 'libero.it', 'outlook.it', 'alice.it', 'aruba.it', 'icloud.com'],
        'ja_JP': ['gmail.com', 'yahoo.co.jp', 'outlook.jp'],
        'ko_KR': ['gmail.com', 'naver.com', 'daum.net', 'hanmail.net'],
        'zh_CN': ['qq.com', '163.com', 'sina.com', '126.com'],
        'ar_SA': ['hotmail.com', 'gmail.com', 'yahoo.com', 'outlook.com'],
        'ru_RU': ['mail.ru', 'yandex.ru', 'gmail.com']
    }

    lista_comune = ["123456","password","123456789","12345","12345678","qwerty",
                    "abc123","111111","123123","admin","letmein",
                    "welcome","monkey","login","iloveyou","sunshine",
                    "princess","dragon","football","baseball","password1",
                    "1234","000000","qwerty123","1q2w3e4r","zaq12wsx","123qwe",
                    "654321","123321","superman","batman","trustno1","hello",
                    "freedom","whatever","qazwsx","ashley","michael",
                    "charlie","andrew","daniel","jessica","hannah","love",
                    "secret","summer","winter","master","shadow",
                    "killer","starwars","pokemon","ninja","welcome123",
                    "passw0rd","admin123","root",
                    "test","user","guest","default","computer"]
    
    locale = random.choice(locales)
    fake = Faker(locale)

    nome = fake.name().replace(" ", "")

    if random.random() < 0.1:
        car = random.choice(['.', '_', '-'])
        luogo = random.randint(1, len(nome) - 1)
        nome = nome[:luogo] + car + nome[luogo:]

    if locale in email_providers:
        provider = random.choice(email_providers[locale])
    else:
        provider = random.choice(email_providers['en_US'])

    email = f"{nome}@{provider}"

    scelta = random.random()

    if scelta < 0.4:
        password = nome.lower() + str(random.randint(39, 99))
    elif scelta < 0.5:
        caratteri = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789'
        password = ''.join(random.choices(caratteri, k=random.randint(6, 12)))
    else:
        prob = random.randint(1, 2)
        if prob == 1:
            password = random.choice(lista_comune)
        else:
            pw_senza_numero = [pwd for pwd in lista_comune if not any(c.isdigit() for c in pwd)]
            password = random.choice(pw_senza_numero)

    print(f"Landscape Credentials  Email: {email} - Password {password}")

# Strati del payload

# LAYER 1 - INITIAL SEED
def layer_1(pwd: bytes) -> bytes:
    salt = struct.pack(">Q", time.time_ns())
    return hashlib.pbkdf2_hmac("sha256", pwd, salt, 20000)

# LAYER 2 - SYSTEM FINGERPRINT HASH
def layer_2(data: bytes) -> bytes:
    fp = (platform.system() + platform.release() + platform.machine()).encode()
    return hashlib.sha256(data + fp).digest()

# LAYER 3 - ENTROPY FUSION
def layer_3(data: bytes) -> bytes:
    noise = secrets.token_bytes(32)
    return hashlib.sha512(data + noise).digest()

# LAYER 4 - TIME WINDOW BLENDING
def layer_4(data: bytes) -> bytes:
    t1 = struct.pack(">Q", time.time_ns())
    t2 = struct.pack(">Q", int(time.time() * 1e6))
    return hashlib.sha256(data + t1 + t2).digest()

# LAYER 5 - BYTE DIFFUSION (SHIFT MIX)
def layer_5(data: bytes) -> bytes:
    return bytes((b << 1 | b >> 7) & 0xFF for b in data)

# LAYER 6 - DOUBLE HASH CASCADE
def layer_6(data: bytes) -> bytes:
    return hashlib.sha256(hashlib.sha512(data).digest()).digest()

# LAYER 7 - XOR STREAM BLEND
def layer_7(data: bytes) -> bytes:
    key = hashlib.sha256(data).digest()
    return bytes(b ^ key[i % len(key)] for i, b in enumerate(data))

# LAYER 8 - MEMORY SALT REINJECTION
def layer_8(data: bytes) -> bytes:
    return hashlib.sha512(data + secrets.token_bytes(64)).digest()

# LAYER 9 - BLOCK PERMUTATION
def layer_9(data: bytes) -> bytes:
    block = 4
    chunks = [data[i:i+block] for i in range(0, len(data), block)]
    chunks.reverse()
    return b"".join(chunks)

# LAYER 10 - STREAM HASH MIX
def layer_10(data: bytes) -> bytes:
    return hashlib.sha256(data + hashlib.sha256(data).digest()).digest()

# LAYER 11 - ENTROPY EXPANSION
def layer_11(data: bytes) -> bytes:
    return hashlib.pbkdf2_hmac("sha512", data, data[:16], 5000)

# LAYER 12 - SELF-FEEDBACK HASH
def layer_12(data: bytes) -> bytes:
    return hashlib.sha512(data + hashlib.sha256(data).digest()).digest()

# LAYER 13 - NONLINEAR MIXING
def layer_13(data: bytes) -> bytes:
    return bytes((a ^ (b + i)) & 0xFF for i, (a, b) in enumerate(zip(data, reversed(data))))

# LAYER 14 - TIME-SALT RESEED
def layer_14(data: bytes) -> bytes:
    return hashlib.sha256(data + str(time.time_ns()).encode() + secrets.token_bytes(16)).digest()

# LAYER 15 - HASH TREE FOLDING
def layer_15(data: bytes) -> bytes:
    mid = len(data) // 2
    left = hashlib.sha256(data[:mid]).digest()
    right = hashlib.sha256(data[mid:]).digest()
    return hashlib.sha512(left + right).digest()

# LAYER 16 - BYTE SHUFFLE SEEDING
def layer_16(data: bytes) -> bytes:
    seed = int.from_bytes(hashlib.sha256(data).digest(), "big")
    lst = list(data)
    for i in range(len(lst)):
        j = (seed + i * 31) % len(lst)
        lst[i], lst[j] = lst[j], lst[i]
    return bytes(lst)

# LAYER 17 - CROSS HASH BLEND
def layer_17(data: bytes) -> bytes:
    h1 = hashlib.sha256(data).digest()
    h2 = hashlib.sha512(data).digest()
    return bytes(a ^ b for a, b in zip(h1, h2[:32]))

# LAYER 18 - BASE64 COMPRESSION
def layer_18(data: bytes) -> bytes:
    return base64.b64encode(data)

# LAYER 19 - FINAL CHAOTIC HASH
def layer_19(data: bytes) -> bytes:
    return hashlib.sha3_512(data + secrets.token_bytes(32)).digest()

# LAYER 20 - OUTPUT FINALE
def layer_20(data: bytes) -> str:
    final = hashlib.sha256(data).hexdigest()
    return final.upper()

# MAIN FUNCTION - COMPONIMENTO PAYLOAD
def final_seed():
    rounds = int(input("Inserisci il numero di rounds: "))

    data = layer_1(sentiero)

    for _ in range(rounds):
        data = layer_2(data)
        data = layer_3(data)
        data = layer_4(data)
        data = layer_5(data)
        data = layer_6(data)
        data = layer_7(data)
        data = layer_8(data)
        data = layer_9(data)
        data = layer_10(data)
        data = layer_11(data)
        data = layer_12(data)
        data = layer_13(data)
        data = layer_14(data)
        data = layer_15(data)
        data = layer_16(data)
        data = layer_17(data)
        data = layer_18(data)
        data = layer_19(data)

    return layer_20(data)


print("Landscape Fingerprint: " + FingerPrintGenerator())
credentials_faker()
print("Landscape Final Payload: " + final_seed())
fine = input()
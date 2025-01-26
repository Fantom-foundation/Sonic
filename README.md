# ⚡ Sonic 

**Sonic is an EVM-compatible blockchain secured by the Lachesis consensus algorithm.**  
Built for **high-speed transactions and decentralized applications**.

<div align="center">

[![GitHub Repo stars](https://img.shields.io/github/stars/Fantom-foundation/Sonic?logo=github&color=yellow)](https://github.com/Fantom-foundation/Sonic/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/Fantom-foundation/Sonic?logo=github&color=blue)](https://github.com/Fantom-foundation/Sonic/network/members)
[![GitHub last commit](https://img.shields.io/github/last-commit/Fantom-foundation/Sonic?logo=git)](https://github.com/Fantom-foundation/Sonic/commits/main)
[![License](https://img.shields.io/github/license/Fantom-foundation/Sonic?logo=open-source-initiative)](https://github.com/Fantom-foundation/Sonic/blob/main/LICENSE)
[![Discord](https://img.shields.io/discord/924442927399313448?logo=discord&color=5865F2)](https://discord.gg/3Ynr2QDSnB)
[![Twitter Follow](https://img.shields.io/twitter/follow/SonicLabs?style=flat&logo=twitter)](https://x.com/SonicLabs)

</div>

---

## 🛠 **Building the Source**

To build **Sonic**, you need **Go (1.21+)** and a **C compiler**.  
Once dependencies are installed, run:

```sh  
make all  
```

The build outputs are the `build/sonicd` and `build/sonictool` executables.

---

## 🗄 **Initializing the Sonic Database**

To join a network, you need a **genesis file**.  
Check [this repository](https://github.com/Fantom-foundation/lachesis_launch) for instructions on obtaining the latest version.

Once you have the **genesis file**, initialize your database with:

```sh 
sonictool --datadir=<target DB path> genesis <path to the genesis file>
```

---

## 🚀 **Running `sonicd`**

Here are some common ways to launch your **Sonic node**.

### 🌍 **Launching a Network Node** (Non-Validator Mode)
Run `sonicd` as a **readonly (non-validator) node** using the genesis file:

```sh  
sonicd --datadir=<DB path>  
```

---

### ⚙️ **Using a Configuration File**
Instead of using multiple CLI flags, pass a **config file**:

```sh  
sonicd --datadir=<DB path> --config /path/to/your/config.toml  
```

To generate a default config file:

```sh  
sonictool --datadir=<DB path> dumpconfig  
```

---

## 🔐 **Running a Validator Node**
To **create a validator private key**, run:

```sh  
sonictool --datadir=<DB path> validator new  
```

To launch a **validator node**, you need to specify your **Validator ID** and **Public Key**:

```sh 
sonicd --datadir=<DB path> --validator.id=YOUR_ID --validator.pubkey=0xYOUR_PUBKEY  
```

📌 **Note:** `sonicd` will prompt for a **password** to decrypt your validator key.  
To use a password file, add:  

```sh
--validator.password /path/to/password-file
```

For details on registering your **validator stake**, see the [Fantom Documentation](https://docs.fantom.foundation).

---

## 🌐 **Network Connectivity**
To improve **network discovery**, specify your **public IP**:

```sh  
sonicd --datadir=<DB path> --nat=extip:1.2.3.4  
```

📌 **Make sure** that TCP/UDP **port 5050** is open in your firewall.

---

## 💬 **Join the Community**
<p align="left">
  <a href="https://t.me/Sonic_English">
    <img src="https://img.shields.io/badge/Telegram-26A5E4?logo=telegram&logoColor=white&style=for-the-badge" alt="Telegram">
  </a>
  <a href="https://discord.gg/3Ynr2QDSnB">
    <img src="https://img.shields.io/badge/Discord-5865F2?logo=discord&logoColor=white&style=for-the-badge" alt="Discord">
  </a>
  <a href="https://x.com/SonicLabs">
    <img src="https://img.shields.io/badge/Twitter-000000?logo=x&logoColor=white&style=for-the-badge" alt="Twitter (X)">
  </a>
</p>

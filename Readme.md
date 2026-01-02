# Circa

Circa is a blockchain-based rotating savings application (Ajo / ROSCA) that allows private groups to pool funds periodically and receive payouts in turn, enforced by smart contracts.

It combines traditional rotating savings with transparent, trust-minimized settlement on-chain, while keeping group metadata and user profiles private off-chain.

---

## What is Ajọ?

Ajọ (also known as ROSCA – Rotating Savings and Credit Association) is a savings system where:

- A fixed group of members agree to contribute a fixed amount periodically
- Each period, the total pooled amount is paid out to one member
- The payout rotates until every member has received the pot once

Circa automates this system using smart contracts to enforce contributions and payouts, removing the need for a trusted organizer.

---

## Core Principles

- **On-chain enforcement**  
  All money movement and payout rules are enforced by smart contracts.

- **Off-chain privacy**  
  Group metadata, invites, and user profiles are managed off-chain to support private groups.

- **Wallet-native identity**  
  Users authenticate using their wallet (sign-in with signature). No passwords.

- **Separation of concerns**  
  - Smart contracts handle correctness and settlement
  - Backend handles indexing, privacy, and UX
  - Frontend handles interaction and presentation

---

## Architecture Overview



### Smart Contracts
- Enforce contribution rules
- Track payments per period
- Release payouts only when fully funded
- Emit events for indexing and auditing

### Backend (Go)
- Indexes on-chain events
- Persists a read-model in Postgres
- Manages:
  - user profiles
  - private groups
  - invite codes
  - access control
- Exposes REST APIs for the frontend

### Frontend (React + Vite)
- Wallet connection and signing
- Group creation and invites
- Contribution and payout actions
- Progress tracking and activity views

---

## Features (MVP)

### Groups
- Create private Ajo groups
- Invite members via invite codes
- Fixed membership per round

### Rounds
- Fixed contribution amount
- Fixed payout order
- One payout per period
- Transparent progress tracking

### Payments
- One contribution per member per period
- Idempotent and verifiable on-chain
- Automatic payout when fully funded

### Profiles
- Wallet-based identity
- Display name and avatar
- No passwords

---

## What Circa Does NOT Do (Yet)

- Handle defaults or late payments
- Apply penalties or interest
- Replace members mid-round
- Provide strong on-chain privacy guarantees
- Support variable contribution amounts

These are intentionally deferred to keep the core system simple and correct.

---

## Tech Stack

### Backend
- Go
- Postgres
- Wallet-based authentication (SIWE-style)

### Smart Contracts
- Solidity
- EVM-compatible chains (local/testnet initially)

### Frontend
- React
- Vite
- Wallet integration (e.g. viem / wagmi)

---

## Repository Structure


# ENZYMES.md — 32 Enzim Restriksi Tipe II

Tabel lengkap 32 enzim restriksi tipe II yang dipakai proyek
`restriction-mapper`. Semua situs pengenalan dan posisi potong di bawah
ini **terverifikasi dari daftar NEB** (*Alphabetized List of
Recognition Specificities*), silang-periksa dengan REBASE,
**diakses 2026-11-20**.

## Cara membaca tabel

- **Situs pengenalan** ditulis 5'→3' pada untai atas, dalam notasi NEB:
  tanda `/` adalah posisi potong untai atas. Contoh: `G/AATTC`
  berarti enzim memotong untai atas setelah basa pertama.
- **Panjang situs** adalah jumlah basa pada situs pengenalan (huruf
  ambiguitas IUPAC dihitung satu basa, misalnya `N` pada `G/ANTC`).
- **Ujung** menyatakan jenis ujung yang dihasilkan, dihitung dengan
  rumus **panjang overhang = |panjang situs − 2 × posisi potong
  atas|**, di mana posisi potong atas adalah jumlah basa sebelum
  tanda `/`:
  - posisi potong < tengah situs → ujung lengket 5' (*5' overhang*);
  - posisi potong > tengah situs → ujung lengket 3' (*3' overhang*);
  - posisi potong tepat di tengah situs → ujung tumpul (*blunt*,
    overhang 0 nt).
- **Organisme sumber** hanya diisi bila penulis yakin dari pengetahuan
  umum yang mapan. Tanda `-` berarti sengaja dikosongkan, bukan berarti
  organisme sumbernya tidak diketahui — penulis tidak ingin mengarang
  dari ingatan.

## Tabel enzim

| Enzim | Situs pengenalan (5'→3') | Panjang situs | Ujung | Organisme sumber |
|-------|--------------------------|---------------|-------|------------------|
| EcoRI | G/AATTC | 6 | 5' overhang 4 nt | *Escherichia coli* RY13 |
| BamHI | G/GATCC | 6 | 5' overhang 4 nt | *Bacillus amyloliquefaciens* H |
| HindIII | A/AGCTT | 6 | 5' overhang 4 nt | *Haemophilus influenzae* Rd |
| XbaI | T/CTAGA | 6 | 5' overhang 4 nt | - |
| SalI | G/TCGAC | 6 | 5' overhang 4 nt | - |
| PstI | CTGCA/G | 6 | 3' overhang 4 nt | - |
| KpnI | GGTAC/C | 6 | 3' overhang 4 nt | - |
| SacI | GAGCT/C | 6 | 3' overhang 4 nt | - |
| SmaI | CCC/GGG | 6 | Blunt | - |
| NcoI | C/CATGG | 6 | 5' overhang 4 nt | - |
| NdeI | CA/TATG | 6 | 5' overhang 2 nt | - |
| XhoI | C/TCGAG | 6 | 5' overhang 4 nt | - |
| SpeI | A/CTAGT | 6 | 5' overhang 4 nt | - |
| NotI | GC/GGCCGC | 8 | 5' overhang 4 nt | - |
| BglII | A/GATCT | 6 | 5' overhang 4 nt | - |
| ClaI | AT/CGAT | 6 | 5' overhang 2 nt | - |
| EcoRV | GAT/ATC | 6 | Blunt | - |
| HpaI | GTT/AAC | 6 | Blunt | - |
| SphI | GCATG/C | 6 | 3' overhang 4 nt | - |
| ApaI | GGGCC/C | 6 | 3' overhang 4 nt | - |
| AccI | GT/MKAC | 6 | 5' overhang 2 nt | - |
| HinfI | G/ANTC | 5 | 5' overhang 3 nt | - |
| TaqI | T/CGA | 4 | 5' overhang 2 nt | - |
| NheI | G/CTAGC | 6 | 5' overhang 4 nt | - |
| AgeI | A/CCGGT | 6 | 5' overhang 4 nt | - |
| MluI | A/CGCGT | 6 | 5' overhang 4 nt | - |
| PvuII | CAG/CTG | 6 | Blunt | - |
| StuI | AGG/CCT | 6 | Blunt | - |
| ScaI | AGT/ACT | 6 | Blunt | - |
| XmaI | C/CCGGG | 6 | 5' overhang 4 nt | - |
| AvaI | C/YCGRG | 6 | 5' overhang 4 nt | - |
| AflII | C/TTAAG | 6 | 5' overhang 4 nt | - |

## Ringkasan distribusi

| Jenis ujung | Jumlah enzim |
|-------------|--------------|
| 5' overhang 4 nt | 16 |
| 5' overhang 2 nt | 4 |
| 5' overhang 3 nt | 1 |
| 3' overhang 4 nt | 5 |
| Blunt | 6 |
| **Total** | **32** |

## Catatan

1. **Semua 32 enzim di atas palindromik.** Situs pengenalan yang sama
   dibaca dari kedua untai 5'→3'. Potongan untai bawah adalah cermin
   dari posisi potong untai atas terhadap titik tengah situs, sehingga
   satu tanda `/` sudah cukup untuk menentukan kedua potongan.
2. **Situs ambigu (kode IUPAC).** Tiga enzim memakai kode ambiguitas:
   AccI `GT/MKAC` (M = A atau C; K = G atau T), HinfI `G/ANTC`
   (N = basa apa pun), dan AvaI `C/YCGRG` (Y = C atau T; R = A atau
   G). Uraian varian konkretnya ada di `RESEARCH.md` Bagian 7.
3. **SmaI dan XmaI mengenali situs yang sama, `CCCGGG`, tetapi
   memotong di tempat berbeda:** SmaI `CCC/GGG` menghasilkan ujung
   tumpul, sedangkan XmaI `C/CCGGG` menghasilkan 5' overhang 4 nt.
4. **Ujung yang kompatibel.** BamHI (`G/GATCC`) dan BglII (`A/GATCT`)
   sama-sama menghasilkan overhang 5' `GATC`, sehingga fragmen yang
   dihasilkannya dapat diligasi satu sama lain. Lihat `RESEARCH.md`
   Bagian 6 untuk pembahasan ujung lengket, ujung tumpul, dan ujung
   yang kompatibel.
5. **Sumber data:** NEB *Alphabetized List of Recognition Sequences*
   dan REBASE, diakses 2026-11-20. Latar belakang ilmiah (palindrom,
   frekuensi pemotongan, tipe IIS, pUC19) ada di `RESEARCH.md`.
6. **Enzim tipe IIS (mis. BsaI) tidak termasuk** dalam tabel ini,
   karena enzim tipe IIS memotong di luar situs pengenalannya dan
   tidak dapat dimodelkan dengan notasi `/` di dalam situs. Lihat
   `RESEARCH.md` Bagian 9.

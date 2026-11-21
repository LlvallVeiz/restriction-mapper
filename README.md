# restriction-mapper

Pemetaan situs restriksi enzim pada sekuens DNA, plus simulasi digesti
(single/double digest, linear/sirkular). Ditulis dalam Go, tanpa
dependensi eksternal.

* **32 enzim restriksi tipe II** dalam database bawaan (EcoRI, BamHI,
  HindIII, … NotI, AccI, HinfI, AvaI) dengan situs pengenalan IUPAC,
  posisi potong 5'/3', dan jenis ujung — semuanya terverifikasi dari
  daftar NEB (lihat `docs/ENZYMES.md`).
* **Scan dua untai** (forward + reverse complement), situs tumpang
  tindih, situs yang membungkus origin pada plasmid sirkular, dan
  daftar enzim yang **tidak** memotong.
* **Simulasi digesti**: single + double digest, DNA linear dan
  sirkular, daftar fragmen + panjangnya (opsional: sekuens fragmen).
* **Statistik**: jumlah situs per enzim vs. jumlah yang diharapkan
  pada DNA acak, peta situs ASCII, output teks/tabel/JSON.
* **CLI** dengan subcommand `scan`, `digest`, `enzymes`, `map` +
  mode batch multi-FASTA.

## Latar belakang

Enzim restriksi tipe II memotong DNA untai ganda pada situs
pengenalan pendek (umumnya 4–8 bp) yang palindromik — situs yang sama
dibaca dari kedua untai. Potongan yang tidak simetris meninggalkan
ujung lengket (*sticky end*, overhang 5' atau 3'), potongan tepat di
tengah meninggalkan ujung tumpul (*blunt end*). Peta situs-situs ini
adalah langkah pertama kloning klasik: memilih enzim yang memotong
tepat sekali di dalam *multiple cloning site* (MCS) sebuah plasmid,
atau yang **tidak** memotong gen yang akan dikloning. Penjelasan
lengkap, termasuk sejarah, kode IUPAC, frekuensi harapan, dan enzim
tipe IIS yang sengaja tidak dimodelkan, ada di
[`docs/RESEARCH.md`](docs/RESEARCH.md).

## Instalasi

Butuh Go 1.22 atau lebih baru. Di VM build proyek ini Go **belum**
terpasang dan dipasang saat build via apt (Ubuntu):

```sh
sudo apt-get update && sudo apt-get install -y golang-go
go version   # go version go1.22.2 linux/amd64
```

Lalu:

```sh
git clone https://github.com/LlvallVeiz/restriction-mapper.git
cd restriction-mapper
go build -o rmap ./cmd/rmap
```

Tidak ada dependensi eksternal (`go.mod` hanya mendeklarasikan modul);
`go test ./...` berjalan tanpa akses internet.

## Penggunaan

```sh
# Daftar 32 enzim dalam database
./rmap enzymes

# Scan semua enzim pada sebuah plasmid (sirkular)
./rmap scan -in examples/input/puc19.fasta -circular

# Hanya enzim tertentu, output JSON
./rmap scan -in examples/input/puc19.fasta -circular \
    -enzymes EcoRI,BamHI,HindIII -format json

# Double digest: fragmen EcoRI + HindIII pada pUC19 sirkular
./rmap digest -in examples/input/puc19.fasta -enzymes EcoRI,HindIII -circular

# Peta ASCII
./rmap map -in examples/input/puc19.fasta -circular -width 100

# Batch: satu file FASTA multi-record, satu laporan per record
./rmap scan -in examples/input/multi.fasta
```

Opsi penting:

| Flag | Arti |
| --- | --- |
| `-in FILE` | File FASTA input (wajib untuk `scan`/`digest`/`map`) |
| `-enzymes A,B` | Nama enzim dipisah koma (default `scan`/`map`: semua 32; `digest`: wajib, maks. 2) |
| `-circular` | Perlakukan sekuens sebagai plasmid sirkular: situs boleh membungkus origin, digest membungkus origin |
| `-format` | `text` (default), `table`, atau `json` |
| `-seq` | Sertakan sekuens tiap fragmen pada output `digest` |
| `-out FILE` | Tulis output ke file, bukan stdout |
| `-width N` | Lebar peta ASCII dalam kolom (default 100) |
| `-id RECORD` | Digest hanya record FASTA dengan ID ini (default: record pertama) |

Posisi memakai koordinat **1-based pada untai atas**. Posisi potong
dilaporkan sebagai "jumlah basa di kiri potongan": EcoRI `G/AATTC`
pada situs yang mulai di basa 396 memotong untai atas setelah basa
396 (potongan di antara basa 396 dan 397).

## Struktur proyek

```
restriction-mapper/
├── cmd/rmap/main.go            CLI (subcommand scan/digest/enzymes/map)
├── internal/
│   ├── dna/                    complement, reverse complement, GC, validasi IUPAC
│   ├── enzyme/                 definisi enzim + database 32 enzim + pencocokan IUPAC
│   ├── fasta/                  parser FASTA multi-record
│   ├── scan/                   pemindaian dua untai, sirkular, statistik
│   ├── digest/                 simulasi digesti linear/sirkular
│   └── report/                 tabel situs, tabel statistik, peta ASCII
├── tests/e2e_test.go           tes end-to-end lewat biner CLI yang sebenarnya
├── docs/
│   ├── RESEARCH.md             latar belakang ilmiah + References
│   └── ENZYMES.md              tabel 32 enzim + sumber
├── examples/
│   ├── input/                  puc19.fasta (nyata, dari NEB), plasmid
│   │                           sintetik, fragmen mirip-lambda, multi.fasta
│   └── output/                 output contoh dari run beneran
└── scripts/
    ├── run-examples.sh         regenerasi examples/output/
    └── verify-puc19.py         verifikasi independen (Python) atas scan pUC19
```

## Contoh hasil

### pUC19 — plasmid kloning nyata (GenBank L09137.2, 2686 bp, sirkular)

Sekuens diunduh dari NEB (`examples/input/puc19.fasta`), diverifikasi
2686 bp, GC 50,63%. Perintah:

```sh
./rmap scan -in examples/input/puc19.fasta -circular
```

Hasil ringkas (output lengkap di
[`examples/output/puc19-scan.txt`](examples/output/puc19-scan.txt)):

```
== pUC19 (2686 bp, GC 50.63%, circular) ==
pUC19: 2686 bp, 27 sites from 18 enzymes, 14 enzymes do not cut
  Pos   Enzyme  Strand  Site      TopCut  BotCut
  396   EcoRI   +       G/AATTC   396     400
  402   SacI    +       GAGCT/C   406     402
  408   KpnI    +       GGTAC/C   412     408
  412   SmaI    +       CCC/GGG   414     414
  412   XmaI    +       C/CCGGG   412     416
  417   BamHI   +       G/GATCC   417     421
  423   XbaI    +       T/CTAGA   423     427
  429   SalI    +       G/TCGAC   429     433
  435   PstI    +       CTGCA/G   439     435
  441   SphI    +       GCATG/C   445     441
  447   HindIII +       A/AGCTT   447     451
  ...
Non-cutters (14): AflII, AgeI, ApaI, BglII, ClaI, EcoRV, HpaI, MluI,
  NcoI, NheI, NotI, SpeI, StuI, XhoI
```

Klaster MCS muncul persis di posisi **396–451** dalam orientasi
L09137.2 — EcoRI, SacI, KpnI, SmaI/XmaI, BamHI, XbaI, SalI, PstI,
SphI, HindIII — masing-masing tepat satu situs, sesuai peta NEB.
Enzim 4-cutter yang sering memotong (HinfI 6 situs, TaqI 4 situs)
mendekati harapan statistiknya (10,49 untuk sekuens 2686 bp), dan
NotI (8-cutter, harapan 0,04) memang tidak memotong.

Double digest EcoRI + HindIII pada pUC19 sirkular:

```sh
./rmap digest -in examples/input/puc19.fasta -enzymes EcoRI,HindIII -circular
```

```
Digest of pUC19 (2686 bp, circular) with EcoRI + HindIII
  EcoRI: 1 site(s)
  HindIII: 1 site(s)
Cut positions (bases left of cut): [396 447]
Fragments (2):
  #    Start   End     Length
  1     397     447     51 bp
  2     448     396     2635 bp
Total: 2686 bp (check: 2686 bp)
```

Fragmen 51 bp adalah potongan MCS di antara kedua situs — cocok
dengan selisih posisi potong 447 − 396 = 51. Satu enzim saja pada
plasmid sirkular hanya melinearkan plasmid (1 fragmen 2686 bp);
lihat `examples/output/puc19-digest-ecori.txt`.

### Plasmid sintetik pSYN3200 (3200 bp, sirkular, dirancang)

Dibangkitkan dengan seed tetap (`random.seed(20261120)`): latar acak
+ 29 situs MCS yang ditanam pada posisi yang diketahui (101, 201,
301, …, 2901 — lihat `scripts/` untuk generatornya; posisi tanam
tercatat di file contoh). Scan menemukan **91 situs dari 30 enzim**
(2 enzim tidak memotong: HpaI, NotI); sisanya di luar 29 situs tanam
adalah situs kebetulan pada latar acak, sesuai harapan statistik.
Double digest EcoRI + BamHI (potong di 101 dan 201):

```
Cut positions (bases left of cut): [101 201]
Fragments (2):
  1     102     201     100 bp
  2     202     101     3100 bp
Total: 3200 bp (check: 3200 bp)
```

### Fragmen mirip-lambda (5000 bp, linear, sintetik)

**112 situs dari 29 enzim**; EcoRI 3 situs (500, 1500, 4500),
BamHI 5 situs, HindIII 1 situs (3500) — termasuk situs kebetulan
pada latar acak. Output lengkap di
`examples/output/lambda-scan.txt`.

### Batch multi-FASTA

`examples/input/multi.fasta` berisi 3 record pendek; satu perintah
`scan` melaporkan ketiganya (`examples/output/batch-scan.txt`):

```
frag-a: 24 bp, 2 sites from 2 enzymes, 30 enzymes do not cut
frag-b: 18 bp, 3 sites from 2 enzymes, 30 enzymes do not cut
frag-c: 12 bp, 0 sites from 0 enzymes, 32 enzymes do not cut
```

## Keputusan desain

* **Satu nilai potong per enzim.** Semua enzim di database bersifat
  palindromik dan memotong di dalam situsnya, sehingga potongan
  untai bawah adalah cermin potongan untai atas
  (`BottomCut = len(Site) − TopCut`). Ini membuat model data kecil
  dan mudah diaudit; konsekuensinya enzim tipe IIS (mis. BsaI,
  memotong di luar situs) tidak masuk database — lihat
  `docs/RESEARCH.md` bagian 9.
* **Pencocokan IUPAC langsung, bukan regex.** Setiap posisi diperiksa
  basa per basa terhadap himpunan IUPAC (`M` = A/C, `N` = apa pun).
  Lebih lambat dari regex untuk sekuens raksasa, tetapi perilakunya
  eksplisit dan mudah diuji — termasuk situs degenerat seperti AccI
  `GTMKAC` (4 varian konkret) dan HinfI `G/ANTC`.
* **Deduplikasi situs palindromik.** Situs palindromik ditemukan baik
  pada pemindaian untai atas maupun reverse complement; keduanya
  adalah satu situs fisik yang sama, jadi dilaporkan sekali (untai
  `+`). Situs non-palindromik (tidak ada di database bawaan, tetapi
  didukung kodenya) tetap dilaporkan dengan untai `-` yang benar.
* **Posisi potong sebagai "basa di kiri potongan".** Satu bilangan
  bulat per potongan membuat aritmetika fragmen (linear/sirkular,
  termasuk potongan di origin) menjadi pengurangan sederhana, dan
  total panjang fragmen selalu sama dengan panjang molekul — dicek
  di setiap output digest.
* **Tanpa dependensi.** Hanya standard library Go; `go test ./...`
  dan semua contoh berjalan offline.

## Tes

```sh
go test ./...   # 45 tes, semuanya lolos
go vet ./...    # bersih
```

Cakupan tes (semua ekspektasi dihitung tangan, bukan disalin dari
kode):

* `dna`: reverse complement (termasuk palindrom EcoRI), komplemen
  kode ambiguitas, GC dengan kontribusi fraksional, validasi.
* `enzyme`: 32 enzim valid + palindromik; overhang tiap enzim
  (EcoRI 5'/4 nt `AATT`, PstI 3'/4 nt `TGCA`, SmaI blunt, AccI
  5'/2 nt degenerat, NotI 5'/4 nt); pencocokan IUPAC (AccI cocok
  ke-4 variannya, AvaI menolak `CTCGTG` karena R = A/G); frekuensi
  harapan (EcoRI 1/4096, TaqI 1/256, NotI 1/65536).
* `scan`: situs EcoRI/BamHI pada posisi yang diketahui tangan
  (5 dan 15) beserta posisi potongnya; tidak ada hitungan ganda
  palindrom; situs tumpang tindih; situs membungkus origin plasmid
  sirkular (ditemukan sirkular, tidak ditemukan linear); situs
  hanya-untai-bawah pada enzim non-palindromik; enzim degenerat.
* `digest`: linear single/double digest, sirkular 0/1/2/3 potong
  (termasuk potongan tepat di origin), fragmen membungkus origin,
  sekuens fragmen.
* `fasta`: multi-record, CRLF, huruf kecil, baris 200 kb, error
  (header kosong, sekuens kosong, basa invalid).
* `tests/e2e_test.go`: membangun biner CLI yang sebenarnya dan
  menjalankan `enzymes`, `scan`, `digest`, dan batch pada contoh —
  termasuk angka headline pUC19 (27 situs, EcoRI 396, fragmen
  double digest 51 + 2635 bp) dan plasmid sintetik (100 + 3100 bp).

**Verifikasi independen.** `scripts/verify-puc19.py` menghitung ulang
semua situs pUC19 dengan implementasi Python yang terpisah (tidak
berbagi kode dengan Go) dan membandingkan jumlah + posisi per enzim:

```
python total sites: 27, go total sites: 27
EcoRI+HindIII top cuts: [396, 447]
EcoRI+HindIII circular fragments (sorted desc): [2635, 51]
VERIFY OK: python and go agree on all 32 enzymes.
```

## Keterbatasan & pengembangan

* **Tidak ada metilasi.** Dam/Dcm/CpG metilasi (mis. DpnI butuh DNA
  termetilasi Dam; XbaI terblokir metilasi Dam yang tumpang tindih)
  tidak dimodelkan; semua situs dianggap dapat dipotong.
* **Hanya enzim tipe II pemotong-dalam-situs.** Enzim tipe IIS
  (BsaI, BsmBI, FokI, …) dan nicking enzyme tidak didukung — butuh
  model potong `(a/b)` di luar situs.
* **Tidak ada suhu/buffer.** Kondisi reaksi (suhu inkubasi, buffer,
  star activity) di luar cakupan; alat ini memetakan situs, bukan
  mensimulasikan kondisi lab.
* **Performa.** Pemindaian O(panjang × enzim); cukup untuk plasmid
  dan genom bakteri kecil, belum dioptimalkan (mis. Aho–Corasick)
  untuk genom eukariot utuh.
* **Pengembangan yang wajar:** impor database enzim dari file
  (format REBASE), dukungan tipe IIS, anotasi metilasi, output
  GenBank/SnapGene, dan peta melingkar (SVG).

## Referensi

* REBASE — The Restriction Enzyme Database —
  <https://rebase.neb.com/rebase/rebase.html>
* NEB, Alphabetized List of Recognition Sequences —
  <https://www.neb.com/en-gb/tools-and-resources/selection-charts/alphabetized-list-of-recognition-specificities>
* NEB, pUC19 sequence (FASTA) —
  <https://www.neb.com/en/-/media/nebus/page-images/tools-and-resources/interactive-tools/dna-sequences-and-maps/text-documents/puc19fsa.txt>
* NCBI GenBank L09137.2 (pUC19c) —
  <https://www.ncbi.nlm.nih.gov/nuccore/L09137.2>
* Yanisch-Perron, Vieira & Messing (1985), Gene 33:103–119 —
  Improved M13 phage cloning vectors (pUC19).
* Roberts RJ et al. (2023), Nucleic Acids Res 51(D1):D629–D637 —
  REBASE 2023. DOI: <https://doi.org/10.1093/nar/gkac1036>
* NEBcutter — <https://nc2.neb.com/NEBcutter2/>

Daftar lengkap + pembahasan ada di
[`docs/RESEARCH.md`](docs/RESEARCH.md).

## Lisensi

MIT — lihat [`LICENSE`](LICENSE).

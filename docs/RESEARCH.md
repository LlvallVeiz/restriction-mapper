# RESEARCH.md — Restriction Mapper

Dokumen latar belakang ilmiah untuk proyek Go `restriction-mapper`: pemetaan situs pemotongan enzim restriksi tipe II pada sekuens DNA.

**Status verifikasi — baca dulu sebelum memakai angka di dokumen ini:**

- **Terverifikasi dari daftar NEB, diakses 2026-11-20:** situs pengenalan
dan posisi potong 32 enzim yang dipakai proyek ini (lihat `ENZYMES.md`). Sumbernya: NEB *Alphabetized List of Recognition Specificities* dan REBASE.
- **Pengetahuan umum yang mapan:** sejarah enzim restriksi, Hadiah Nobel
1978, konsep palindrom DNA, ujung lengket vs. tumpul, kode ambiguitas IUPAC, rumus frekuensi pemotongan yang diharapkan, dan identitas plasmid pUC19. Fakta-fakta ini sudah lama menjadi konsensus dalam buku teks biologi molekuler dan tidak diklaim sebagai hasil pengukuran baru proyek ini.
- **Tidak ada hasil eksperimen baru di dokumen ini.** Tidak ada angka
pemotongan yang diukur di laboratorium, tidak ada sekuens baru, dan tidak ada sitasi yang dikarang. Semua rujukan ada di bagian References, dengan URL persis seperti yang tercantum di sana.

## Daftar isi

Konvensi notasi · 1. Apa itu enzim restriksi? · 2. Sejarah singkat
3. Sistem restriksi–modifikasi · 4. Enzim restriksi tipe II
5. Palindrom DNA · 6. Ujung lengket vs. ujung tumpul
7. Kode ambiguitas IUPAC · 8. Frekuensi pemotongan yang diharapkan
9. Enzim tipe IIS · 10. Aplikasi · 11. REBASE · 12. pUC19
13. Keterbatasan model · Glosarium singkat · 14. References

## Konvensi notasi yang dipakai di dokumen ini

Agar tidak ada salah tafsir, berikut konvensi yang dipakai konsisten di `RESEARCH.md` dan `ENZYMES.md`:

- Sekuens DNA selalu ditulis 5'→3' pada untai atas, kecuali dinyatakan lain.
- Tanda `/` di dalam situs pengenalan adalah posisi potong untai atas, mengikuti notasi NEB (misalnya `G/AATTC` berarti potong setelah basa pertama).
- `bp` berarti pasangan basa; `nt` berarti nukleotida (dipakai untuk panjang overhang untai tunggal).
- `L` adalah panjang situs pengenalan dalam bp; `p` adalah posisi potong untai atas, yaitu jumlah basa sebelum tanda `/`.
- `n` adalah panjang situs pengenalan saat dipakai dalam rumus frekuensi 4^(−n).
- Huruf R, Y, S, W, K, M, B, D, H, V, N adalah kode ambiguitas IUPAC (Bagian 7), bukan basa tunggal.
- Nama enzim ditulis dengan huruf pertama kapital dan akhiran Romawi (EcoRI, BamHI, HindIII), mengikuti konvensi penamaan REBASE/NEB.
- "Terverifikasi" berarti situs dan posisi potongnya cocok dengan daftar NEB yang diakses 2026-11-20; "pengetahuan umum yang mapan" berarti fakta dari buku teks dan literatur standar, bukan hasil pengukuran proyek ini.
- Semua perhitungan harapan (frekuensi, jumlah situs yang diharapkan) diberi label sebagai harapan statistik, bukan hasil yang diamati.

## 1. Apa itu enzim restriksi?

Enzim restriksi — disebut juga **endonuklease restriksi** — adalah enzim yang memotong DNA untai ganda pada (atau di dekat) situs pengenalan (*recognition site*) yang sekuensnya spesifik. Setiap enzim mengenali satu sekuens pendek, umumnya 4–8 pasangan basa (bp), lalu memotong ikatan fosfodiester pada kedua untai DNA.

Tiga sifatnya yang membuatnya menjadi perkakas dasar biologi molekuler adalah **spesifisitas** (hanya memotong sekuens yang dikenali), **reprodusibilitas** (sekuens yang sama selalu dipotong di posisi yang sama), dan **kompatibilitas ujung** (fragmen yang dipotong enzim yang sama, atau enzim berbeda dengan ujung yang kompatibel, dapat disambung kembali dengan ligase). Proyek `restriction-mapper` memodelkan persis perilaku yang dapat diprediksi ini: diberikan sekuens DNA dan satu set enzim, ia menghitung di mana saja setiap enzim akan memotong — tanpa memodelkan kinetika, efisiensi, atau kondisi reaksi.

## 2. Sejarah singkat

**2.1. Pengamatan awal.** Pada 1950-an dan 1960-an, para peneliti mengamati bahwa bakteriofag yang tumbuh pada satu galur bakteri sering gagal menginfeksi galur lain — seolah pertumbuhan fag itu "direstriksi" oleh inang barunya. DNA fag masuk ke dalam sel, lalu dihancurkan di dalamnya. **Werner Arber** merumuskan penjelasan yang kemudian terbukti benar: bakteri membedakan DNA sendiri dari DNA asing. DNA sendiri dilindungi oleh penanda kimia (metilasi), sedangkan DNA asing yang tidak bermetilasi dengan benar dikenali dan dipotong.

**2.2. Isolasi enzim pertama.** **Hamilton O. Smith** mengisolasi dan mengkarakterisasi enzim restriksi tipe II pertama, **HindII**, dari *Haemophilus influenzae* pada tahun 1970. HindII memotong DNA di tengah situs pengenalannya yang spesifik — perilaku yang bersih dan dapat diprediksi, berbeda dari enzim tipe I yang memotong pada jarak acak dari situs pengenalannya. Karakter inilah yang membuat enzim tipe II langsung berguna di laboratorium.

**2.3. Pemetaan DNA dan Hadiah Nobel 1978.** **Daniel Nathans** menunjukkan kegunaan praktis yang menentukan: ia memakai enzim restriksi untuk memetakan genom virus SV40 — memotong DNA virus menjadi fragmen, lalu menyusun peta genetiknya. Atas pekerjaan ini, **Hadiah Nobel Fisiologi atau Kedokteran 1978** diberikan bersama kepada **Werner Arber, Hamilton O. Smith, dan Daniel Nathans** "untuk penemuan enzim restriksi dan penerapannya pada masalah genetika molekuler". Tiga kontribusi itu saling melengkapi: Arber dengan gagasannya, Smith dengan enzimnya, Nathans dengan penerapannya.

**2.4. Setelah 1978.** Ribuan enzim restriksi ditemukan dari berbagai bakteri dan arke, dikarakterisasi situs pengenalannya, dan dikatalogkan. Perusahaan seperti New England Biolabs (NEB) memurnikan dan menjualnya sebagai reagen standar, sementara REBASE (lihat Bagian 11) menjadi katalog acuan bagi seluruh komunitas. Kloning molekuler modern — dan dengan itu sebagian besar bioteknologi rekombinan — berdiri di atas fondasi ini.

> Uraian sejarah di atas adalah pengetahuan umum yang mapan dari
> literatur biologi molekuler. Tahun dan nama yang disebutkan adalah
> yang sudah lama tercatat dalam sumber-sumber standar.

## 3. Sistem restriksi–modifikasi: pertahanan bakteri

Enzim restriksi adalah senjata bakteri, bukan alat yang dirancang untuk manusia. Setiap sistem restriksi–modifikasi (R–M) terdiri dari dua komponen berpasangan: **endonuklease restriksi** yang memotong DNA yang tidak terlindungi, dan **metiltransferase** yang menambahkan gugus metil pada situs pengenalan yang sama di DNA bakteri sendiri, sehingga enzim restriksi pasangannya tidak memotongnya. Logikanya adalah diskriminasi diri vs. bukan-diri: DNA bakteri sendiri bermetilasi → tidak dipotong; DNA fag yang masuk tidak bermetilasi dengan pola yang benar → dipotong menjadi fragmen dan infeksi digagalkan.

Nama enzim mencerminkan asal-usulnya: huruf pertama genus, dua huruf pertama spesies, diikuti nomor galur dan nomor urut. Contoh: **EcoRI** dari *Escherichia coli* galur RY13, **BamHI** dari *Bacillus amyloliquefaciens* galur H, **HindIII** dari *Haemophilus influenzae* galur Rd. Tiga asal-usul inilah yang dicantumkan di `ENZYMES.md`; untuk enzim lain, kolom organisme sumber sengaja dikosongkan daripada diisi dari ingatan yang tidak yakin.

## 4. Enzim restriksi tipe II

Enzim restriksi dibagi menjadi beberapa tipe (I, II, III, dan seterusnya) berdasarkan cara kerjanya. Yang relevan untuk proyek ini — dan untuk hampir semua pekerjaan kloning — adalah **tipe II**. Ciri-cirinya: **memotong di dalam situs pengenalannya sendiri** (atau tepat di tepinya), bukan pada jarak acak darinya, sehingga posisi potongnya dapat diprediksi dari sekuens saja; **tidak membutuhkan ATP** untuk memotong; **situs pengenalan umumnya palindromik** (Bagian 5); umumnya membutuhkan ion Mg²⁺; dan pada kebanyakan enzim tipe II klasik, satu polipeptida membawa baik fungsi pengenalan maupun pemotongan.

Semua 32 enzim dalam `ENZYMES.md` adalah enzim tipe II dengan situs pengenalan palindromik. Itulah yang membuat proyek ini dapat memodelkannya dengan satu aturan seragam: cari situs pengenalan pada kedua untai, lalu potong pada posisi yang sudah ditentukan. Enzim tipe IIS — subkelompok tipe II yang memotong **di luar** situs pengenalannya — sengaja **tidak** dimodelkan; lihat Bagian 9.

## 5. Palindrom DNA

Situs pengenalan enzim restriksi tipe II hampir selalu **palindromik**. Yang dimaksud palindrom dalam konteks DNA berbeda dari palindrom kata seperti "katak": sebuah sekuens DNA disebut palindromik bila **sekuens untai atas yang dibaca 5'→3' sama dengan sekuens untai bawah yang dibaca 5'→3'** — dengan kata lain, sekuens tersebut sama dengan komplemen balik (*reverse complement*) dirinya sendiri.

Contoh — EcoRI:

```text
5'– G A A T T C –3'
3'– C T T A A G –5'
```

Untai atas 5'→3' adalah `GAATTC`. Untai bawah, bila dibaca 5'→3' (dari kanan ke kiri pada gambar di atas), juga `GAATTC`. Karena itu enzim yang sama mengenali situs tersebut dari kedua arah, dan pemotongan pada untai bawah adalah cermin dari pemotongan pada untai atas. Contoh lain: BamHI `GGATCC`, HindIII `AAGCTT`, SmaI `CCCGGG` — ketiganya palindromik dengan definisi yang sama.

**Mengapa palindrom?** Enzim restriksi tipe II umumnya bekerja sebagai **homodimer**: dua subunit protein yang identik, masing-masing mengenali separuh situs pengenalan pada satu untai. Situs yang palindromik memungkinkan kedua subunit yang identik itu mengikat kedua untai dengan cara yang simetris. Simetri ini juga yang membuat pemotongan pada kedua untai saling berhubungan secara geometris: posisi potong untai bawah ditentukan sepenuhnya oleh posisi potong untai atas.

**Konsekuensi untuk pemodelan.** Karena semua enzim dalam proyek ini palindromik, pencarian situs cukup dilakukan pada satu untai; posisi potong untai bawah adalah cermin posisi potong untai atas, sehingga tidak perlu data tambahan; dan notasi NEB dengan satu tanda `/` (misalnya `G/AATTC`) sudah cukup untuk menentukan kedua potongan. Tidak semua situs pengenalan bersifat palindromik sempurna — ada enzim tipe II dengan situs asimetris — tetapi tidak satu pun termasuk dalam 32 enzim proyek ini.

## 6. Ujung lengket vs. ujung tumpul

Ketika enzim memotong kedua untai DNA, kedua potongan bisa terjadi pada posisi yang sama atau pada posisi yang bergeser. Dari perbedaan inilah muncul dua jenis ujung.

**Ujung lengket (*sticky end* / *cohesive end*)** terjadi bila kedua untai dipotong pada posisi yang **berbeda**, sehingga salah satu untai menjulur sebagai untai tunggal sepanjang beberapa nukleotida. Ujung yang menjulur ini dapat berpasangan kembali (*anneal*) dengan ujung komplementer dari fragmen lain melalui ikatan hidrogen — "lengket" secara harfiah pada tingkat molekuler. Ada dua macam: **ujung lengket 5' (*5' overhang*)**, bila posisi potong untai atas berada **sebelum titik tengah** situs pengenalan; dan **ujung lengket 3' (*3' overhang*)**, bila posisi potong untai atas berada **setelah titik tengah** situs pengenalan.

Contoh — EcoRI (`G/AATTC`), ujung lengket 5' sepanjang 4 nt:

```text
5'– G         A A T T C –3'
3'– C T T A A         G –5'
```

Setelah dipotong, kedua fragmen masing-masing membawa ujung tunggal `AATT` yang saling komplementer. Contoh — PstI (`CTGCA/G`), ujung lengket 3' sepanjang 4 nt, di mana untai bawah yang menjulur:

```text
5'– C T G C A         G –3'
3'– G         A C G T C –5'
```

**Ujung tumpul (*blunt end*)** terjadi bila kedua untai dipotong **tepat di posisi yang sama** — di tengah situs pengenalan — sehingga tidak ada untai tunggal yang menjulur. Contoh — SmaI (`CCC/GGG`):

```text
5'– C C C         G G G –3'
3'– G G G         C C C –5'
```

Ujung tumpul dapat diligasi ke ujung tumpul mana pun, tetapi dengan efisiensi yang jauh lebih rendah daripada ujung lengket: tidak ada ikatan hidrogen pendahuluan yang menahan kedua ujung tetap bersama selama ligase bekerja, dan tidak ada spesifisitas arah.

**Aturan penghitungan yang dipakai proyek ini.** Untuk situs sepanjang *L* dan posisi potong untai atas *p* (jumlah basa sebelum tanda `/` dalam notasi NEB): **panjang overhang = |L − 2 × p|**. Bila *p < L/2* → ujung lengket 5' sepanjang |L − 2p| nt; bila *p > L/2* → ujung lengket 3' sepanjang |L − 2p| nt; bila *p = L/2* → ujung tumpul (overhang 0 nt). Contoh: EcoRI `G/AATTC`, L = 6, p = 1 → 4 nt, 5' overhang; PstI `CTGCA/G`, L = 6, p = 5 → 4 nt, 3' overhang; SmaI `CCC/GGG`, L = 6, p = 3 → 0 nt, tumpul; HinfI `G/ANTC`, L = 5, p = 1 → 3 nt, 5' overhang; TaqI `T/CGA`, L = 4, p = 1 → 2 nt, 5' overhang. Rumus yang sama menghasilkan kolom "Ujung" di `ENZYMES.md` untuk seluruh 32 enzim.

**Distribusi jenis ujung pada 32 enzim proyek ini.** Menghitung kolom "Ujung" di `ENZYMES.md` dengan rumus di atas menghasilkan ringkasan berikut — angka-angka ini adalah hasil penghitungan dari tabel itu sendiri, bukan data baru:

| Jenis ujung | Panjang overhang | Jumlah enzim |
|-------------|------------------|--------------|
| 5' overhang | 4 nt | 16 |
| 5' overhang | 2 nt | 4 |
| 5' overhang | 3 nt | 1 |
| 3' overhang | 4 nt | 5 |
| Tumpul (blunt) | 0 nt | 6 |
| **Total** | | **32** |

Pola yang terlihat: ujung lengket 5' sepanjang 4 nt adalah yang paling umum (setengah dari himpunan ini), yang sejalan dengan dominannya enzim 6-cutter dengan potongan pada posisi 1 atau 5 dalam praktik kloning. Ujung tumpul hanya 6 dari 32 enzim, dan satu-satunya overhang 3 nt yang ganjil panjangnya — 3 nt pada HinfI — berasal dari situsnya yang panjangnya ganjil (5 bp), karena hanya situs dengan panjang ganjil yang dapat menghasilkan overhang ganjil dari rumus |L − 2p|.

Satu catatan menarik: **SmaI dan XmaI mengenali situs yang sama, `CCCGGG`, tetapi memotong di tempat berbeda** (`CCC/GGG` tumpul vs. `C/CCGGG` dengan 5' overhang 4 nt). Keduanya menghasilkan ujung yang kompatibel untuk ligasi, tetapi fragmen yang dihasilkan tidak identik — perbedaan yang hanya terlihat bila posisi potong dicatat dengan benar, bukan hanya situsnya. Lebih umum lagi, **ujung yang kompatibel** (*compatible ends*) dapat berasal dari dua enzim berbeda: BamHI (`G/GATCC`) dan BglII (`A/GATCT`) sama-sama menghasilkan overhang 5' `GATC`, sehingga fragmen BamHI dapat diligasi ke fragmen BglII — walaupun situs hibrid yang terbentuk umumnya tidak lagi dikenali oleh kedua enzim tersebut. Perencana kloning yang baik selalu memeriksa hal ini, karena ujung lengket memberi spesifisitas dan efisiensi ligasi, sedangkan ujung tumpul bersifat universal tetapi tidak terarah dan kurang efisien.

## 7. Kode ambiguitas nukleotida IUPAC

Tidak semua posisi dalam situs pengenalan bersifat tetap. Beberapa enzim mengenali sekuens yang pada satu atau beberapa posisi dapat berisi lebih dari satu basa. Untuk menulis sekuens semacam itu, dipakai **kode ambiguitas nukleotida IUPAC** — satu huruf untuk setiap himpunan basa yang mungkin. Tabel lengkap (15 kode; standar IUPAC, lihat References):

| Kode | Basa yang diwakili | Arti / mnemonik |
|------|--------------------|-----------------|
| A | A | Adenin |
| C | C | Sitosin |
| G | G | Guanin |
| T | T | Timin (U dalam RNA) |
| R | A atau G | **R** = pu**R**in |
| Y | C atau T | p**Y**rimidin |
| S | G atau C | **S** = ikatan *strong* (3 ikatan H) |
| W | A atau T | **W** = ikatan *weak* (2 ikatan H) |
| K | G atau T | **K** = *keto* |
| M | A atau C | **M** = a**M**ino |
| B | C, G, atau T | bukan A |
| D | A, G, atau T | bukan C |
| H | A, C, atau T | bukan G |
| V | A, C, atau G | bukan T |
| N | A, C, G, atau T | basa apa pun (**N** = *any*) |

Cara membacanya: huruf tunggal A, C, G, T berarti posisi itu tetap; **R, Y, S, W, K, M** masing-masing berarti dua kemungkinan basa; **B, D, H, V** masing-masing berarti tiga kemungkinan basa ("semua kecuali satu"); **N** berarti empat kemungkinan basa — posisi itu tidak memberi informasi apa pun.

Contoh dari 32 enzim proyek ini: **AccI** `GT/MKAC` — `M` berarti posisi ketiga dapat berupa A atau C, dan `K` berarti posisi keempat dapat berupa G atau T; situs ini mencakup empat varian konkret (`GTAGAC`, `GTATAC`, `GTCGAC`, `GTCTAC`), dan salah satu variannya, `GTCGAC`, adalah situs SalI. **HinfI** `G/ANTC` — `N` berarti posisi ketiga dapat berupa basa apa pun (empat varian: `GAATC`, `GACTC`, `GAGTC`, `GATTC`). **AvaI** `C/YCGRG` — `Y` (C atau T) pada posisi kedua dan `R` (A atau G) pada posisi kelima; empat varian konkret.

Kode ambiguitas memperbesar jumlah situs yang harus diperiksa. Cara yang paling jujur adalah menguraikan situs ambigu menjadi semua varian konkretnya, lalu mencari masing-masing varian — atau memakai pola pencocokan yang memahami kode IUPAC secara langsung. Proyek ini mendokumentasikan situs ambigu persis seperti yang tercantum di NEB (`GT/MKAC`, `G/ANTC`, `C/YCGRG`), bukan dengan menggantinya dengan salah satu variannya. Perilaku pencocokan yang tepat menjadi tanggung jawab implementasi dan harus diuji terhadap uraian varian di atas.

## 8. Frekuensi pemotongan yang diharapkan

Seberapa sering sebuah enzim memotong DNA acak? Jawabannya mengikuti aturan peluang yang sederhana — dan penting untuk dipahami sebagai **harapan statistik, bukan janji**. Untuk sekuens acak dengan keempat basa yang sama-sama mungkin (probabilitas 1/4 per posisi), peluang bahwa suatu posisi tertentu memulai situs pengenalan sepanjang *n* adalah **P = (1/4)^n = 4^(−n)**, sehingga jarak yang diharapkan antar situs — atau panjang fragmen rata-rata — adalah **4^n bp**.

| Panjang situs (n) | Peluang per posisi | Jarak yang diharapkan |
|-------------------|--------------------|-----------------------|
| 4 (*4-cutter*) | 1/256 | ~256 bp |
| 5 | 1/1024 | ~1024 bp |
| 6 (*6-cutter*) | 1/4096 | ~4096 bp |
| 8 (*8-cutter*) | 1/65536 | ~65536 bp |

Inilah asal angka-angka yang dipakai di mana-mana: **4-cutter** (misalnya TaqI, `TCGA`) memotong rata-rata setiap **~256 bp** — sering, menghasilkan banyak fragmen kecil; **6-cutter** (misalnya EcoRI, BamHI, HindIII) memotong rata-rata setiap **~4096 bp** — jarang, menghasilkan fragmen berukuran kilobasa yang cocok untuk kloning; **8-cutter** (misalnya NotI, `GCGGCCGC`) memotong rata-rata setiap **~65536 bp** — sangat jarang, berguna untuk memetakan genom besar.

**Harapan pada plasmid pUC19 (2686 bp).** Dengan memakai panjang pUC19 (Bagian 12) sebagai contoh perhitungan — **ini adalah perhitungan harapan, bukan jumlah pemotongan yang diamati**: 4-cutter, 2686 / 256 ≈ **10,5 situs yang diharapkan**; 6-cutter, 2686 / 4096 ≈ **0,66 situs yang diharapkan** — kebanyakan 6-cutter tidak memotong pUC19 sama sekali, dan yang memotong biasanya hanya sekali; 8-cutter, 2686 / 65536 ≈ **0,04 situs yang diharapkan** — NotI hampir pasti tidak memotong pUC19. Inilah sebabnya plasmid kloning dirancang dengan *multiple cloning site* (MCS): situs-situs 6-cutter sengaja ditempatkan di sana, karena secara acak hampir tidak ada 6-cutter yang memotong plasmid berukuran ~2,7 kb. Tanpa MCS, plasmid itu hampir tidak dapat dipotong secara berguna.

**Mengapa angka ini hanya harapan.** Rumus 4^(−n) bertumpu pada tiga asumsi yang **tidak pernah persis benar** pada DNA nyata: (1) keempat basa sama-sama mungkin — genom nyata memiliki bias kandungan GC, jadi situs yang kaya GC (misalnya SmaI `CCCGGG`) lebih jarang pada genom yang miskin GC, dan sebaliknya; (2) posisi saling bebas — DNA nyata memiliki bias dinukleotida (misalnya CpG yang jarang pada genom mamalia) dan sekuens berulang; (3) situs ambigu dihitung setara — untuk situs dengan kode IUPAC (misalnya HinfI `GANTC`), peluang efektifnya lebih besar daripada 4^(−n): `N` dapat berupa basa apa pun, jadi peluangnya adalah 4^(−4) = 1/256, sama seperti 4-cutter, bukan 4^(−5) = 1/1024; menghitung frekuensi situs ambigu memerlukan penjumlahan atas semua varian konkretnya. Karena itu angka-angka di atas dipakai sebagai **panduan kasar untuk memilih enzim**, bukan sebagai prediksi jumlah potongan pada sekuens tertentu. Jumlah potongan yang sebenarnya hanya diketahui dengan memindai sekuens itu sendiri — persis itulah yang dilakukan `restriction-mapper`.

## 9. Enzim tipe IIS — keterbatasan yang tidak dimodelkan

Tidak semua enzim tipe II memotong **di dalam** situs pengenalannya. Ada subkelompok, **tipe IIS**, yang mengenali satu sekuens tetapi memotong pada jarak tetap **di luar** sekuens itu. Contoh kanoniknya adalah **BsaI: GGTCTC(1/5)** — BsaI mengenali `GGTCTC`, lalu memotong untai atas 1 nt setelah situs dan untai bawah 5 nt setelah situs, di luar situs pengenalan itu sendiri. Konsekuensinya penting: setelah pemotongan, situs pengenalan BsaI **tidak lagi berada pada fragmen produk** — ia terbuang bersama fragmen kecil di antara dua situs. Karena itu BsaI (dan enzim tipe IIS lainnya) dapat dipakai untuk perakitan tanpa bekas (*scarless assembly*) — prinsip yang dipakai dalam **Golden Gate cloning** (Bagian 10).

**Mengapa tipe IIS tidak dimodelkan di proyek ini.** Model `restriction-mapper` mengasumsikan situs pengenalan palindromik dan posisi potong ditentukan oleh satu tanda `/` di dalam situs, dengan potongan untai bawah sebagai cerminnya. Kedua asumsi itu **gagal** untuk enzim tipe IIS: situs tipe IIS umumnya **tidak palindromik** (BsaI: `GGTCTC` — komplemen baliknya `GAGACC`); posisi potongnya berada di luar situs dan dinyatakan dengan notasi berbeda, `(1/5)`, yang tidak dapat direpresentasikan dengan satu `/` di dalam situs; dan pemotongan untai bawah bukan cermin sederhana dari potongan untai atas terhadap titik tengah situs. Memasukkan tipe IIS membutuhkan model data yang berbeda (situs asimetris, offset pemotongan luar, dan arah enzim). Proyek ini sengaja **tidak** memodelkannya, dan BsaI tidak termasuk dalam 32 enzim di `ENZYMES.md`. Keputusan ini didokumentasikan di sini agar tidak ada yang salah mengira cakupan proyek lebih luas dari kenyataannya.

## 10. Aplikasi

**10.1. Kloning molekuler.** Alur kloning restriksi klasik: potong vektor plasmid dan fragmen insert dengan enzim yang sama (atau enzim dengan ujung yang kompatibel); campurkan keduanya, ujung lengket yang komplementer berpasangan; ligase DNA menyambung tulang punggung fosfodiester; transformasikan ke bakteri dan seleksi koloni yang membawa plasmid rekombinan. Pemilihan enzim menentukan segalanya: enzim harus memotong vektor tepat sekali (di MCS), tidak memotong insert di tempat yang tidak diinginkan, dan menghasilkan ujung yang kompatibel. Tanpa peta restriksi yang akurat, langkah pertama ini sudah gagal — itulah gunanya perangkat pemetaan seperti proyek ini.

**10.2. RFLP — *Restriction Fragment Length Polymorphism*.** Sebelum sekuensing DNA menjadi murah, RFLP adalah salah satu teknik utama untuk mendeteksi variasi genetik: DNA genomik dipotong dengan enzim restriksi; bila sebuah mutasi (misalnya SNP) menciptakan atau menghilangkan situs pengenalan, pola panjang fragmen berubah; pola fragmen yang berbeda antar individu menjadi penanda genetik (*marker*) yang dapat diikuti dalam keluarga atau populasi. RFLP dipakai untuk pemetaan genetik, diagnosis penyakit (contoh klasik: anemia sel sabit mengubah situs pengenalan), dan sidik jari DNA forensik pada masa awalnya. Prinsipnya persis sama dengan yang dimodelkan proyek ini — perbedaan satu basa dalam situs pengenalan mengubah hasil pemotongan — sehingga ketepatan pencocokan situs adalah segalanya.

**10.3. Golden Gate cloning.** Golden Gate adalah metode perakitan multi-fragmen dalam satu reaksi yang memanfaatkan enzim tipe IIS, terutama BsaI dan BsmBI: setiap fragmen dirancang dengan situs BsaI di kedua ujungnya, mengapit overhang 4 nt yang berbeda-beda; BsaI memotong di luar situs pengenalannya (Bagian 9), sehingga situs pengenalan terbuang dan overhang yang tersisa menentukan urutan perakitan; karena produk akhir tidak lagi mengandung situs BsaI, reaksi hanya berjalan satu arah — menuju produk yang benar. Golden Gate menunjukkan mengapa batasan model proyek ini penting untuk dinyatakan dengan jujur: metode yang paling modern dalam perakitan DNA justru bergantung pada enzim yang **tidak** dimodelkan di sini. Pengguna yang membutuhkan Golden Gate harus memakai perangkat yang mendukung tipe IIS.

**10.4. Aplikasi lain.** **Pemetaan restriksi (*restriction mapping*):** menentukan posisi situs pengenalan pada DNA yang belum disekuensing, dengan memotongnya memakai beberapa enzim dan mengukur panjang fragmen — inilah yang dilakukan Nathans pada SV40. **Subkloning dan verifikasi konstruk:** memotong plasmid hasil kloning dengan enzim diagnostik dan memeriksa pola fragmen pada gel, cara cepat memastikan konstruk benar sebelum disekuensing. **Pustaka genomik dan cDNA:** pemotongan parsial dengan 4-cutter menghasilkan fragmen berukuran acak yang cocok untuk kloning ke dalam vektor.

## 11. REBASE sebagai basis data acuan

**REBASE** (*Restriction Enzyme dataBASE*) adalah katalog acuan untuk enzim restriksi dan metiltransferase. REBASE mencatat, untuk setiap enzim yang diketahui: nama, situs pengenalan, posisi potong, organisme sumber, metilase pasangannya, sensitivitas terhadap metilasi, serta ketersediaan komersialnya. Publikasi acuan yang dipakai proyek ini: Roberts RJ et al. (2023), "REBASE — a database for DNA restriction and modification: enzymes, genes and genomes", *Nucleic Acids Research* 51 (D1): D629–D637, DOI 10.1093/nar/gkac1036 (sitasi jurnal; URL DOI tercantum di References).

Data enzim dalam proyek ini (situs + posisi potong, 32 enzim di `ENZYMES.md`) diverifikasi terhadap **NEB Alphabetized List of Recognition Specificities** dan silang-periksa dengan REBASE, diakses 2026-11-20. Bila terjadi perbedaan antara daftar NEB dan REBASE pada suatu enzim, REBASE adalah otoritas taksonomi dan penamaan, sedangkan data situs/posisi potong yang dipakai proyek adalah yang tercantum di NEB pada tanggal akses tersebut. REBASE juga menjadi pengingat tentang skala bidang ini: jumlah enzim restriksi yang diketahui mencapai ribuan, dengan ratusan situs pengenalan yang berbeda. 32 enzim proyek ini hanyalah himpunan awal yang umum dipakai di laboratorium kloning — bukan katalog lengkap.

## 12. pUC19 sebagai contoh nyata

**pUC19** adalah plasmid kloning standar dan contoh nyata yang paling tepat untuk menguji pemeta restriksi: cukup kecil untuk dipahami seluruhnya, cukup umum sehingga hasilnya dapat diperiksa oleh siapa pun. Fakta-fakta dasar (pengetahuan umum yang mapan; rujukan di References): **panjang 2686 bp**; **akses GenBank L09137.2** (pUC19c); **publikasi** Yanisch-Perron, Vieira, dan Messing (1985), "Improved M13 phage cloning vectors and host strains: nucleotide sequences of the M13mp18 and pUC19 vectors", *Gene* 33:103–119; **sekuens FASTA** tersedia di NEB (lihat References).

Isi plasmid ini menjelaskan mengapa ia begitu berguna: **gen resistensi ampisilin (*bla*)** untuk seleksi; **asal replikasi (*ori*)** untuk replikasi dalam jumlah salinan tinggi di *E. coli*; **gen *lacZα*** untuk skrining biru-putih — penyisipan fragmen di MCS memutus gen ini, sehingga koloni rekombinan berwarna putih dan koloni tanpa insert berwarna biru pada medium X-gal; dan **multiple cloning site (MCS)** — sekelompok situs pengenalan 6-cutter yang berdekatan di dalam *lacZα*, termasuk situs EcoRI, BamHI, HindIII, XbaI, SalI, PstI, KpnI, SacI, SmaI, dan lain-lain yang termasuk dalam 32 enzim proyek ini. Hubungan pUC19 dengan Bagian 8 sudah dijelaskan di sana: tanpa MCS, hampir tidak ada 6-cutter yang akan memotong plasmid 2686 bp. MCS adalah jawaban desain terhadap fakta statistik itu.

Untuk pengujian proyek ini, pUC19 adalah sekuens uji yang ideal: sekuens FASTA resminya tersedia, panjangnya diketahui (2686 bp), dan hasil pemetaannya dapat diverifikasi secara independen dengan perangkat seperti NEBcutter (lihat References). Proyek ini **tidak** mengklaim jumlah potongan tertentu pada pUC19 di dokumen ini — jumlah itu adalah keluaran yang harus dihitung oleh program dari sekuensnya, bukan angka yang ditulis dari ingatan.

## 13. Keterbatasan model proyek ini

Agar tidak ada yang salah memakai perangkat ini, batasan berikut dinyatakan eksplisit: (1) **hanya enzim tipe II palindromik** — enzim tipe IIS (misalnya BsaI) tidak dimodelkan (Bagian 9); (2) **hanya 32 enzim** — himpunan awal, bukan katalog lengkap; ribuan enzim lain tercatat di REBASE; (3) **pencocokan sekuens saja** — model ini tidak mempertimbangkan metilasi DNA (yang dapat memblokir pemotongan pada beberapa enzim), kondisi buffer, suhu inkubasi, aktivitas bintang (*star activity* — pemotongan di situs yang mirip tetapi tidak identik dalam kondisi non-optimal), atau efisiensi pemotongan parsial; (4) **frekuensi yang diharapkan adalah harapan statistik**, bukan prediksi jumlah potongan pada sekuens tertentu (Bagian 8); (5) **organisme sumber hanya dicantumkan bila yakin** — di `ENZYMES.md`, hanya EcoRI, BamHI, dan HindIII yang diberi organisme sumber; sisanya bertanda `-`, bukan karena organisme itu tidak diketahui, melainkan karena penulis dokumen ini tidak ingin mengarangnya dari ingatan; (6) **tidak ada data eksperimen baru** — semua angka di dokumen ini adalah hasil perhitungan dari definisi (rumus overhang, rumus 4^(−n)) atau fakta mapan yang dirujuk, bukan hasil pengukuran.

## Glosarium singkat

- **Bakteriofag (fag):** virus yang menginfeksi bakteri; ancaman yang dihadapi sistem restriksi–modifikasi.
- **Blunt end (ujung tumpul):** ujung DNA hasil pemotongan di mana kedua untai terpotong pada posisi yang sama, tanpa overhang.
- **Cohesive end:** sinonim ujung lengket; ujung DNA dengan overhang untai tunggal yang dapat berpasangan dengan ujung komplementer.
- **Komplemen balik (*reverse complement*):** sekuens yang diperoleh dengan membalik urutan basa dan mengganti setiap basa dengan komplemennya (A↔T, G↔C).
- **Ligase DNA:** enzim yang menyambung kembali tulang punggung fosfodiester antara dua ujung DNA yang berdekatan.
- **MCS (*multiple cloning site*):** daerah pendek pada plasmid yang memuat banyak situs pengenalan enzim restriksi yang berdekatan; disebut juga polylinker.
- **Metilasi:** penambahan gugus metil pada basa DNA; dalam sistem R–M, metilasi menandai DNA sendiri agar tidak dipotong.
- **Metiltransferase:** enzim yang menambahkan gugus metil pada situs pengenalan tertentu.
- **Overhang:** untai tunggal DNA yang menjulur melewati ujung untai pasangannya setelah pemotongan bergeser.
- **Palindrom DNA:** sekuens yang sama dengan komplemen baliknya sendiri (Bagian 5).
- **RFLP:** variasi panjang fragmen restriksi antar individu, dipakai sebagai penanda genetik.
- **Situs pengenalan (*recognition site*):** sekuens DNA spesifik yang dikenali dan dipotong oleh enzim restriksi.
- **Star activity:** pemotongan pada situs yang mirip tetapi tidak identik dengan situs pengenalan, terjadi dalam kondisi reaksi non-optimal.
- **Transformasi:** proses memasukkan DNA plasmid ke dalam sel bakteri.
- **Ujung kompatibel (*compatible ends*):** ujung dari enzim berbeda yang overhang-nya identik sehingga dapat diligasi satu sama lain.
- **Ujung lengket (*sticky end*):** ujung DNA dengan overhang untai tunggal hasil pemotongan bergeser.

## References

Daftar ini berisi **hanya** rujukan yang diberikan untuk proyek ini. Tidak ada URL lain yang ditambahkan, dan tidak ada sitasi yang dikarang.

- **REBASE** — https://rebase.neb.com/rebase/rebase.html
- **NEB Alphabetized List of Recognition Sequences** —
https://www.neb.com/en-gb/tools-and-resources/selection-charts/alphabetized-list-of-recognition-specificities
- **NEB pUC19 sequence (FASTA)** —
https://www.neb.com/en/-/media/nebus/page-images/tools-and-resources/interactive-tools/dna-sequences-and-maps/text-documents/puc19fsa.txt
- **NCBI GenBank L09137.2 (pUC19c)** —
https://www.ncbi.nlm.nih.gov/nuccore/L09137.2
- **Wikipedia pUC19** — https://en.wikipedia.org/wiki/pUC19
- **IUPAC nucleotide codes (NCBI)** —
https://www.ncbi.nlm.nih.gov/sra/docs/submitfile/
- **Roberts RJ et al. REBASE 2023, *Nucleic Acids Research*** —
Roberts RJ et al. "REBASE — a database for DNA restriction and modification: enzymes, genes and genomes." *Nucleic Acids Research* (2023) 51 (D1): D629–D637. https://doi.org/10.1093/nar/gkac1036
- **NEBcutter** — https://nc2.neb.com/NEBcutter2/
- **Cornish-Bowden A (1985), IUPAC codes** — sitasi jurnal saja:
Cornish-Bowden A. "Nomenclature for incompletely specified bases in nucleic acid sequences: recommendations 1984." *Nucleic Acids Research* (1985) 13 (9): 3021–3030.

> Catatan akhir: sitasi Cornish-Bowden 1985 dicantumkan sebagai sitasi
> jurnal saja, tanpa URL, sesuai instruksi penugasan. Sitasi Roberts
> et al. 2023 dicantumkan dengan DOI seperti yang diberikan dalam
> penugasan. Untuk semua rujukan web di atas, tanggal akses yang
> dinyatakan untuk data enzim adalah 2026-11-20.

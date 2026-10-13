# Every app on Earth is now open source... whether it likes it or not

---

**[0:00]** Last  week,  a  developer  got  so  tired  of

**[0:02]** Adobe's  predatory  business  model  that  he

**[0:04]** decided  to  rebuild  seven  of  its  Creative

**[0:06]** Cloud  apps  with  Rust,  like  Photoshop,

**[0:08]** Illustrator,  and  Premiere.  It  then  open

**[0:10]** sourced  them  and  gave  them  away  for

**[0:12]** free.  The  project  is  called  Artcraft,

**[0:14]** and  the  developer's  name  is  Brandon

**[0:16]** Thomas,  who  after  being  hit  with  Adobe's

**[0:18]** cancellation  fee  more  than  once,  opened

**[0:20]** up  Claude,  then  reverse  engineered

**[0:22]** Adobe's  30  plus  year  codebase  from

**[0:24]** scratch,  and  the  apps  even  run  on  Linux,

**[0:26]** a  platform  Adobe  is  never  officially

**[0:28]** supported.  What's  crazy  though  is  that

**[0:30]** you,  the  human  being  watching  this

**[0:31]** video,  can  do  the  same  thing  thanks  to

**[0:33]** an  unrelated  project  called  Reverse

**[0:35]** Engineer  Anything.  It's  been  going  viral

**[0:37]** on  the  internet  lately  because  it  can

**[0:39]** easily  decompile  any  native  binaries  you

**[0:41]** can  get  your  hands  on.  We're  talking

**[0:43]** stripped  ELF  PE  and  MCCO  executables,

**[0:45]** JVM  bite  code  with  obuscated  constant

**[0:47]** pools,  ECMA  335  compliant  assemblies,

**[0:50]** bare  metal  firmware  with  architecture

**[0:52]** specific  instruction  sets,  and  even  OCI

**[0:54]** container  layers  entangled  in

**[0:56]** declarative  infrastructure  manifests.  In

**[0:58]** layman's  terms,  that  means  virtually

**[1:00]** every  app  in  the  world  is  now  open

**[1:01]** source,  whether  it  wants  to  be  or  not.

**[1:03]** In  today's  video,  we'll  attempt  to

**[1:05]** reverse  engineer  anything  and  everything

**[1:07]** and  find  out  if  this  is  the  final  nail

**[1:08]** in  the  coffin  for  software  as  we  know

**[1:10]** it.  It  is  October  12th,  2026,  and  you're

**[1:13]** watching  the  code  report.  So,  is  it

**[1:14]** really  true  that  one  person  plus  Claude

**[1:17]** made  the  entire  Adobe  suite  obsolete

**[1:19]** overnight?  Well,  on  the  surface,  yes.  As

**[1:21]** you  can  see  on  the  website,  you  can

**[1:23]** download  a  Rustbased  version  of  any

**[1:24]** Adobe  program  for  free  right  now.  That's

**[1:27]** awesome.  But  after  trying  it  out  myself,

**[1:29]** I  realized  there's  one  small  catch.  The

**[1:31]** software  is  a  nearly  unusable  pile  of

**[1:33]** vibecoded  horseshit.  Now,  I  begrudgingly

**[1:35]** pay  Adobe  at  $69  tribute  every  month  to

**[1:38]** use  Premiere  to  edit  these  videos  for

**[1:39]** you  guys.  And  I  was  thrilled  to  finally

**[1:41]** see  an  open-  source  alternative.  And

**[1:42]** let's  just  ignore  the  fact  that  Da  Vinci

**[1:44]** Resolve  has  been  an  open  source

**[1:45]** alternative  for  over  20  years  now.  And

**[1:47]** when  I  first  opened  FilmCraft,  I  was

**[1:49]** thrilled  to  see  that  the  UI  was  nearly  a

**[1:51]** onetoone  clone  of  Premiere  so  I  could

**[1:53]** get  to  work  immediately  without  any

**[1:54]** learning  curve.  My  plan  was  to  edit  this

**[1:56]** video  with  it,  but  the  cracks  became

**[1:58]** visible  immediately.  Like  after  dragging

**[2:00]** a  clip  of  some  lame  stock  video  onto  the

**[2:02]** timeline  and  hitting  play,  there  was  a

**[2:04]** massive  lag  in  the  audio  bars  and  just

**[2:06]** overall  slowness  and  laggginess  despite

**[2:08]** being  viodated  in  Rust,  which  I've  been

**[2:10]** told  is  a  blazingly  fast  programming

**[2:12]** language.  So  then  I  tried  implementing

**[2:14]** some  transitions  between  clips  and  they

**[2:16]** were  just  nowhere  near  the  quality  of

**[2:18]** Premiere,  which  you  can  clearly  see  with

**[2:19]** this  glitch  transition  side  by  side.  And

**[2:21]** for  someone  like  me  who's  a  once-  in  a

**[2:23]** generation  cinematic  visionary  who  many

**[2:25]** have  called  the  Stanley  Kubri  of

**[2:27]** sponsored  chip  posts  that  the  sad

**[2:28]** reality  is  that  filmcraft  is  just  not

**[2:30]** good  enough  yet.  But  I  sincerely  hope  it

**[2:32]** gets  there  someday  because  Adobe.  But

**[2:35]** then  that  got  me  thinking  what  if

**[2:36]** instead  of  trying  to  use  off-brand

**[2:38]** Premiere  I  just  reverse  engineer  it

**[2:39]** myself  and  make  it  simpler  because

**[2:41]** Premiere  has  hundreds  of  features  that  I

**[2:43]** don't  use  that  just  get  in  the  way.

**[2:45]** Luckily  there's  a  new  project  in  town

**[2:46]** that  can  help  us  do  that  called  reverse

**[2:48]** engineer  Anything.  Adobe  Premiere  is

**[2:50]** primarily  written  in  C  and  C++  and  crazy

**[2:53]** smart  Asians  on  H1B  visas  have  been

**[2:55]** writing  its  code  for  decades  now.  Its

**[2:57]** actual  source  code  is  top  secret  Adobe

**[2:59]** intellectual  property.  But  in  order  to

**[3:01]** actually  run  it  on  your  machine,  they

**[3:03]** need  to  compile  that  code  into  a  native

**[3:05]** binary  of  ones  and  zeros.  This  stuff  is

**[3:07]** almost  unreadable  to  humans,  but  with  a

**[3:09]** lot  of  painstaking  work,  it  is

**[3:11]** technically  possible  to  map  these  ones

**[3:13]** and  zeros  back  into  something  that

**[3:14]** resembles  human  readable  source  code.

**[3:16]** and  the  REA  project  can  do  that  for

**[3:18]** people  who  can't  even  count  to  three  in

**[3:20]** binary.  First,  we  set  it  up  by  running

**[3:22]** it  with  Node.js  and  the  MPX  command.  And

**[3:25]** then  because  I'm  a  loser  on  Microsoft

**[3:26]** Windows,  we'll  configure  REA's  GEDra

**[3:29]** integration,  which  is  actually  the  tool

**[3:31]** that  deserves  the  most  credit  here  and

**[3:32]** can  do  a  readonly  analysis  of  native

**[3:34]** .exees  and  DLS.  Then  from  there,  we  can

**[3:37]** just  use  it  like  any  other  AI  agent,  and

**[3:39]** tell  it  to  inspect  the  locally  installed

**[3:41]** version  of  Adobe  Premiere.  In  this  case,

**[3:43]** I'll  have  it  inspect  Adobe's  timeline

**[3:44]** integration  and  try  to  figure  out  how  it

**[3:46]** does  this  one  specific  feature.  A  while

**[3:48]** later,  we  get  back  a  bunch  of  sealike

**[3:50]** pseudo  code  along  with  assembly,

**[3:52]** strings,  symbols,  calls,  and  references,

**[3:54]** which  to  an  obsolete  HTML  programmer  is

**[3:57]** all  pretty  useless.  However,  we  can  now

**[3:59]** take  this  garbage  code  and  pass  it  off

**[4:01]** to  our  favorite  coding  agent  to

**[4:02]** reimplement  the  exact  same  feature  in

**[4:04]** Rust.  And  thanks  to  reverse  engineering,

**[4:06]** I'll  eventually  be  able  to  build  my  own

**[4:08]** version  of  Premiere  Light  that's  not

**[4:10]** only  free,  but  only  has  the  actual

**[4:12]** features  I  want  to  use.  But  as  a

**[4:13]** CheapKate  developer,  I'm  not  stopping

**[4:15]** there.  I  think  it's  crazy  that  any  real

**[4:17]** Vibe  developer  in  2026  is  still  paying

**[4:20]** money  for  developer  tools  when  that

**[4:21]** money  should  be  spent  burning  tokens.

**[4:23]** For  example,  let's  take  a  look  at  an

**[4:25]** ancient  technology  called  Supabase,

**[4:27]** which  is  basically  a  wrapper  on  top  of

**[4:29]** AWS  that  provides  a  UI  for  your

**[4:31]** database.  Even  though  it's  mostly  open-

**[4:33]** source,  it  has  a  ton  of  features  I  don't

**[4:35]** want.  So,  what  I'm  doing  instead  is

**[4:37]** using  the  networking  feature  in  Rea  to

**[4:39]** analyze  API  interactions  and  UI  changes,

**[4:42]** which  has  allowed  me  to  build  a  better

**[4:43]** backend  that  I'm  calling  Studbase,  which

**[4:46]** is  not  only  perfectly  tailored  to  the

**[4:48]** needs  of  Horse  Tinder  and  Donker,  but  is

**[4:50]** also  orders  of  magnitude  cheaper  and

**[4:52]** faster  according  to  this  benchmark  I

**[4:53]** just  made  up.  Trust  me,  bro,  you  can

**[4:55]** trust  these  Trust  Me  Bro  benchmarks.  And

**[4:57]** another  thing  you  can  trust  is  Apify,

**[4:59]** the  sponsor  of  today's  video.  Their  MCP

**[5:02]** server  lets  you  connect  your  agent  to

**[5:04]** thousands  of  tools  called  actors.  They

**[5:06]** can  do  things  like  scrape  live  data  from

**[5:08]** YouTube  or  monitor  your  competitor's

**[5:10]** prices  on  Amazon.  But  right  now,  I'm

**[5:12]** using  it  to  track  the  most  liked  posts

**[5:14]** on  the  donkeys  of  Instagram  hashtag  is

**[5:16]** so  I  can  invite  the  most  desirable  asses

**[5:18]** to  sign  up  on  Donker.  It  only  takes  one

**[5:20]** line  to  connect  to  the  server  and  I

**[5:22]** could  explain  what  I  wanted  it  to  do  in

**[5:24]** natural  language.  From  there,  my  agent

**[5:26]** started  using  the  relevant  Appify  tools

**[5:28]** to  return  to  the  structured  data  I

**[5:29]** needed  in  order  to  help  these  gay

**[5:31]** donkeys  find  love.  But  the  most

**[5:33]** impressive  part  was  how  few  tokens  it

**[5:34]** took.  But  because  Appify's  MCP  server

**[5:37]** only  exposes  a  handful  of  core  tools,  it

**[5:39]** then  finds  the  best  one  for  your  job  and

**[5:42]** runs  it  in  the  background  to  avoid  token

**[5:44]** bloat.  You  get  $25  in  free  credits  to

**[5:46]** try  it  out  by  using  the  code  fire  ship

**[5:48]** 25  at  the  link  below.  This  has  been  the

**[5:50]** code  report.  Thanks  for  watching  and  I

**[5:52]** will  see  you  in  the  next  one.

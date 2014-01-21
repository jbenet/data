#!/usr/bin/env python

import os
import re

# constants
readme_file = 'tar.README.md'
semver_regx = r'^[0-9]+\.[0-9]+\.[0-9]+$' # lacks pre-releases/builds
valid_archs = [
  'darwin_amd64',
  'darwin_386',
  'linux_amd64',
  'linux_386',
  'windows_amd64',
  'windows_386',
]


def check(cond, msg):
  if not cond:
    print 'Error:', msg
    exit(-1)

def write_readme(output, arch, version):
  with open(output, 'w') as out:
    with open('../%s' % readme_file) as inp:
      txt = inp.read()
      txt = txt % {'arch': arch, 'version': version}
      out.write(txt)


def make_archive(arch, vers):
  if arch not in valid_archs:
    print "Error: arch '%s' not supported" % arch
    return -1

  if not re.match(semver_regx, vers):
    print "Error: version '%s' is not like X.X.X" % vers
    return -1

  if not os.path.exists('%s/data' % arch):
    print "Error: binary '%s/data' not found" % arch
    return -1

  # move into arch dir
  os.chdir(arch)

  # setup directory
  dir = 'data-v%s-%s' % (vers, arch)
  os.system('mkdir -p %s' % dir)

  # write files
  os.system('cp data %s/data' % dir)
  write_readme('%s/README.md' % dir, arch, vers)

  # tar
  tar = '%s.tar.gz' % dir
  os.system('tar czf %s %s' % (tar, dir))

  # move into place
  os.chdir('..')
  os.system('mkdir -p archives')
  os.system('mv %s/%s archives/%s' % (arch, tar, tar))
  os.system('rm -rf %s/%s' % (arch, dir))

  print 'packaged archives/%s' % tar
  return dir


def main():
  import sys
  if '-h' in sys.argv or len(sys.argv) < 3:
    print 'Usage: %s <arch> <version>' % sys.argv[0]
    print 'Prepares the release archive for a given architecture.'
    exit(0 if '-h' in sys.argv else -1)

  arch = sys.argv[1]
  vers = sys.argv[2]

  archs = valid_archs if arch == 'all' else [arch]

  for arch in archs:
    make_archive(arch, vers)


if __name__ == '__main__':
  main()

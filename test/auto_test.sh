echo "///////////////////////////////////////////"

tgcom-cobra -d true -f prova1.go -l 3-5
tgcom-cobra -d true -f prova2.sh -l 2-4

tgcom-cobra -d true -f prova1.go:3-5,prova2.sh:2-4

echo "///////////////////////////////////////////"

tgcom-cobra -f prova1.go -l 3-5 -a comment
tgcom-cobra -f prova2.sh -l 2-4 -a comment

echo "file uno:"
cat prova1.go

echo ""

echo "file due:"
cat prova2.sh

echo ""

tgcom-cobra -f prova1.go -l 3-5 -a uncomment
tgcom-cobra -f prova2.sh -l 2-4 -a uncomment

echo "file uno:"
cat prova1.go

echo ""

echo "file due:"
cat prova2.sh





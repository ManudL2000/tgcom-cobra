echo "///////////////////////////////////////////"

ciaoo -d true -f prova1.go -l 3-5
ciaoo -d true -f prova2.sh -l 2-4

ciaoo -d true -f prova1.go:3-5,prova2.sh:2-4

echo "///////////////////////////////////////////"

ciaoo -f prova1.go -l 3-5 -a comment
ciaoo -f prova2.sh -l 2-4 -a comment

echo "file uno:"
cat prova1.go

echo ""

echo "file due:"
cat prova2.sh

echo ""

ciaoo -f prova1.go -l 3-5 -a uncomment
ciaoo -f prova2.sh -l 2-4 -a uncomment

echo "file uno:"
cat prova1.go

echo ""

echo "file due:"
cat prova2.sh




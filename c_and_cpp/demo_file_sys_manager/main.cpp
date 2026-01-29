#include "myfilesystemmanager.h"

#include <QApplication>

int main(int argc, char *argv[])
{
	QApplication a(argc, argv);
	QString startPath;
	if (argc > 1)
		startPath = QString(argv[1]);

	MyFileSystemManager mf;
	if (!startPath.isEmpty())
		mf.setStartdemo_fs_Image(startPath);

	mf.show();
	return a.exec();
}

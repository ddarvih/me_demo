#include "myfilesystemmanager.h"

#include <QApplication>
#include <QFileDialog>
#include <QFileInfo>
#include <QHBoxLayout>
#include <QKeyEvent>
#include <QMessageBox>
#include <QString>

MyFileSystemManager::MyFileSystemManager(QWidget* parent) : QMainWindow(parent), isLeftPannelActive(true)
{
	QWidget* w = new QWidget(this);

	QHBoxLayout* pannelsPlace = new QHBoxLayout(w);
	lp = new LeftPannel(this);
	rp = new Rightdemo_fs_Pannel(this);
	pannelsPlace->addWidget(lp);
	pannelsPlace->addWidget(rp);
	w->setLayout(pannelsPlace);
	setCentralWidget(w);

	connect(
		rp,
		&Rightdemo_fs_Pannel::matchRightAsActive,
		this,
		[this]()
		{
			if (isLeftPannelActive)
			{
				rp->setFocus();
			}
			isLeftPannelActive = false;
			updPannelsVisualState();
		});

	connect(
		lp,
		&LeftPannel::matchLeftAsActive,
		this,
		[this]()
		{
			if (!isLeftPannelActive)
				lp->setFocus();
			isLeftPannelActive = true;
			updPannelsVisualState();
		});

	setToolBar();
	lp->setFocus();
	resize(1200, 800);

	updPannelsVisualState();
}

void MyFileSystemManager::setToolBar()
{
	toolBar = addToolBar("Tools");
	toolBar->setMovable(false);
	toolBar->setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Fixed);

	toolBar->addSeparator();

	QWidget* leftSpace = new QWidget(this);
	leftSpace->setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Preferred);
	toolBar->addWidget(leftSpace);

	toolBar->setStyleSheet(
		"QToolBar {"
		"    background: palette(window);"
		"    border: 1px solid palette(mid);"
		"    spacing: 6px;"
		"}"
		"QToolButton {"
		"    background: palette(button);"
		"    border: 1px solid palette(dark);"
		"    padding: 4px 10px;"
		"    margin: 2px;"
		"}");

	addToolBar(Qt::BottomToolBarArea, toolBar);

	QAction* actAbout = toolBar->addAction("[F1] About");
	actAbout->setShortcut(Qt::Key_F1);
	connect(actAbout, &QAction::triggered, this, &MyFileSystemManager::actAboutRequest);

	QAction* actMountF2 = toolBar->addAction("[F2] Mount");
	actMountF2->setShortcut(Qt::Key_F2);
	connect(actMountF2, &QAction::triggered, this, &MyFileSystemManager::mountF2);

	QAction* actMountF4 = toolBar->addAction("[F4] Mount...");
	actMountF4->setShortcut(Qt::Key_F4);
	connect(actMountF4, &QAction::triggered, this, &MyFileSystemManager::mountF4Dialog);

	QAction* actCopy = toolBar->addAction("[F5] Copy");
	actCopy->setShortcut(Qt::Key_F5);
	connect(actCopy, &QAction::triggered, this, &MyFileSystemManager::copyF5);

	QAction* actExit = toolBar->addAction("[F10] Exit");
	actExit->setShortcut(Qt::Key_F10);
	connect(actExit, &QAction::triggered, this, &MyFileSystemManager::exitF10);

	QAction* actChangeFocus = toolBar->addAction("Change active Focus");
	actChangeFocus->setShortcut(Qt::Key_Tab);
	connect(
		actChangeFocus,
		&QAction::triggered,
		this,
		[this]()
		{
			if (isLeftPannelActive)
			{
				rp->setFocus();
				rp->setActive();
			}
			else
			{
				lp->setFocus();
				lp->setActive();
			}
			isLeftPannelActive = !isLeftPannelActive;
			updPannelsVisualState();
		});

	QWidget* rightSpace = new QWidget(this);
	rightSpace->setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Preferred);
	toolBar->addWidget(rightSpace);
}

void MyFileSystemManager::keyPressEvent(QKeyEvent* event)
{
	if (event->key() == Qt::Key_Tab)
	{
		if (isLeftPannelActive)
		{
			rp->setFocus();
			rp->setActive();
		}
		else
		{
			lp->setFocus();
			lp->setActive();
		}
		isLeftPannelActive = !isLeftPannelActive;
		updPannelsVisualState();

		event->accept();
		return;
	}
	QMainWindow::keyPressEvent(event);
}

void MyFileSystemManager::actAboutRequest()
{
	QMessageBox::information(
		this,
		"About",
		"Welcome to My File Manager\n"
		"now it gives you apportunity to work the same way\n"
		"with two almost equal pannels, that have different"
		" realisation\nEnjoy:)");
}

void MyFileSystemManager::mountF2()
{
	if (!isLeftPannelActive)
	{
		QMessageBox::information(this, "Info", "F2 acts with the [Left Pannel: active].");
		return;
	}
	QString path = lp->getFilePath();
	if (path.isEmpty())
	{
		QMessageBox::warning(this, "Fail", "Unable to act with selected demo_fs_ image.");
		return;
	}
	rp->loaddemo_fs_Image(path);
	isLeftPannelActive = false;
	rp->setFocus();
	updPannelsVisualState();
}

void MyFileSystemManager::mountF4Dialog()
{
	QDialog pathGetter(this);
	pathGetter.setWindowTitle("[F4] Mount..");
	pathGetter.resize(400, 150);
	pathGetter.setWindowIconText("Choosing demo_fs_ image");

	QLineEdit* path = new QLineEdit(&pathGetter);
	path->setPlaceholderText("Enter demo_fs_ image to open");
	QPushButton* pressApply = new QPushButton("Apply", &pathGetter);

	QVBoxLayout* getter = new QVBoxLayout(&pathGetter);
	getter->addWidget(path);
	getter->addWidget(pressApply);

	connect(pressApply,
			&QPushButton::clicked,
			[&]()
			{
				QString wantedPath = path->text().trimmed();
				if (!QFile::exists(wantedPath))
				{
					QMessageBox::warning(&pathGetter, "Fail", "Incorrect path");
					pathGetter.accept();
				}
				rp->loaddemo_fs_Image(wantedPath);
				pathGetter.accept();
			});
	pathGetter.exec();

	isLeftPannelActive = false;
	rp->setFocus();
	updPannelsVisualState();
}

void MyFileSystemManager::copyF5()
{
	if (isLeftPannelActive)
	{
		QMessageBox::information(this, "Info", "F5 acts with the [Right Pannel: active].");
		return;
	}

	QVector< FileDescr > filesToCopy = rp->getSelectedFiles();
	if (filesToCopy.isEmpty())
	{
		QMessageBox::warning(this, "Error", "Nothing selected to be to copied");
		return;
	}

	QString destination = lp->getCurDir();

	for (const FileDescr& fd : filesToCopy)
	{
		if (fd.type == "Dir")
		{
			QMessageBox::warning(this, "Error", "[IGNORED] not able to copy DIR");
			continue;
		}

		QString aimPlace = destination + "/" + fd.name;
		if (QFile::exists(aimPlace))
		{
			QMessageBox dublicate(this);
			dublicate.setIcon(QMessageBox::Warning);
			dublicate.setWindowTitle("File exists");
			dublicate.setText(QString("File \"%1\" already exists.\nWant to [skip] this one or [replace]?").arg(fd.name));
			QPushButton* bSkip = dublicate.addButton("Skip", QMessageBox::RejectRole);
			QPushButton* bReplace = dublicate.addButton("Replace", QMessageBox::AcceptRole);
			dublicate.exec();

			if (dublicate.clickedButton() == bSkip)
				continue;
		}

		QByteArray fileContent = rp->getFileContent(fd);
		if (fileContent.isEmpty())
			continue;

		QFile destinationFile(aimPlace);
		if (!destinationFile.open(QIODevice::WriteOnly))
		{
			QMessageBox::StandardButton chooseToFinish =
				QMessageBox::question(this, "[FAIL] to act with destination path", "Want to continue?", QMessageBox::Yes | QMessageBox::No);
			if (chooseToFinish == QMessageBox::No)
			{
				destinationFile.close();
				break;
			}
			continue;
		}
		if (destinationFile.write(fileContent) != fileContent.size())
		{
			QMessageBox::warning(this, "Error", QString("Something went wrong while acting \"%1\" copying ").arg(destination));
		}
		destinationFile.close();
	}
	lp->updWindow();
	QMessageBox::information(this, "Copied", "Finished [F5] Copy");
}

void MyFileSystemManager::exitF10()
{
	QApplication::quit();
}

void MyFileSystemManager::updPannelsVisualState()
{
	QString activeState = "QTreeView { background: palette(base); opacity: 1.0; }";
	QString inactiveState = "QTreeView { background: palette(alternate-base); opacity: 0.8; }";
	lp->findChild< QTreeView* >("listOfFiles")->setStyleSheet(isLeftPannelActive ? activeState : inactiveState);
	rp->findChild< QTreeView* >("listOfFiles")->setStyleSheet(!isLeftPannelActive ? activeState : inactiveState);
	lp->updWindow();
	rp->updWindow();
}

void MyFileSystemManager::setStartdemo_fs_Image(const QString& path)
{
	rp->loaddemo_fs_Image(path);
}

// for later
void MyFileSystemManager::getDirSizeF3() {}
void MyFileSystemManager::recCatalogsCpyF5() {}
void MyFileSystemManager::leftToRightCpyF5() {}
void MyFileSystemManager::replaceElementF6() {}
void MyFileSystemManager::mkDirF7() {}
void MyFileSystemManager::deleteF8() {}

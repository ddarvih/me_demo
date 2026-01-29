#include "rightdemo_fs_pannel.h"

#include "utilspannels.h"

#include <QDebug>
#include <QMessageBox>

Rightdemo_fs_Pannel::Rightdemo_fs_Pannel(QWidget* parent) : QWidget(parent)
{
	demo_fs_FileSystem = new demo_fs_model(this);
	QHeaderView* header;
	QHBoxLayout* pathSetting;
	QVBoxLayout* mainContainer;

	pannelsBase::initPannel(this, &listOfFiles, &curLocationEdit, &backButton, &forwardButton, &header, &pathSetting, &mainContainer);
	listOfFiles->setModel(demo_fs_FileSystem);
	header->setSectionResizeMode(0, QHeaderView::Stretch);
	header->setSectionResizeMode(1, QHeaderView::ResizeToContents);
	header->setSectionResizeMode(2, QHeaderView::ResizeToContents);

	connect(listOfFiles, &QTreeView::doubleClicked, this, &Rightdemo_fs_Pannel::actItemSelected);
	connect(header, &QHeaderView::sectionClicked, this, &Rightdemo_fs_Pannel::actSortingByHeader);
	connect(backButton, &QPushButton::clicked, this, &Rightdemo_fs_Pannel::returnOneBack);
	connect(forwardButton, &QPushButton::clicked, this, &Rightdemo_fs_Pannel::goOneForward);
	setLayout(mainContainer);

	updWindow();
}

void Rightdemo_fs_Pannel::loaddemo_fs_Image(const QString& path)
{
	if (!demo_fs_FileSystem->getdemo_fs_(path))
	{
		QMessageBox::warning(this, "Error", "Failed to load demo_fs_ image. Ensure it is a valid demo_fs_16 or demo_fs_32 image.");
		return;
	}

	pathHistory.backStack.clear();
	pathHistory.forwardStack.clear();

	quint32 mainCluster = demo_fs_FileSystem->isdemo_fs_32TypeModel() ? demo_fs_FileSystem->getRootClusterModel() : 0;
	demo_fs_FileSystem->setCurCluster(mainCluster);
	pathHistory.backStack.push(mainCluster);

	updWindow();
}

QString Rightdemo_fs_Pannel::getCurDir()
{
	return demo_fs_FileSystem->formPath();
}

QString Rightdemo_fs_Pannel::getFilePath()
{
	QModelIndex ind = listOfFiles->currentIndex();

	if (ind.isValid())
		return getCurDir() + "/" + demo_fs_FileSystem->getFileDescript(ind).name;

	return QString();
}

void Rightdemo_fs_Pannel::returnOneBack()
{
	if (pathHistory.backStack.size() <= 1)
		return;
	moveByOne(pathHistory.backStack, pathHistory.forwardStack);
}

void Rightdemo_fs_Pannel::goOneForward()
{
	if (pathHistory.forwardStack.isEmpty())
		return;
	moveByOne(pathHistory.forwardStack, pathHistory.backStack);
}

void Rightdemo_fs_Pannel::moveByOne(QStack< quint32 >& fromStack, QStack< quint32 >& toStack)
{
	toStack.push(demo_fs_FileSystem->getCurCluster());
	quint32 newCluster = fromStack.pop();
	demo_fs_FileSystem->setCurCluster(newCluster);
	listOfFiles->setRootIndex(QModelIndex());
	updWindow();
	actChangedToActive();
}

void Rightdemo_fs_Pannel::actItemSelected(const QModelIndex& ind)
{
	if (!ind.isValid())
		return;

	FileDescr fd = demo_fs_FileSystem->getFileDescript(ind);
	if (!fd.isActiveItem)
		QMessageBox::warning(this, "Warning", "Acting with [DELETED] item");

	if (fd.type == "Dir")
	{
		pathHistory.backStack.push(demo_fs_FileSystem->getCurCluster());
		pathHistory.forwardStack.clear();
		demo_fs_FileSystem->setCurCluster(fd.clusterStart);
		listOfFiles->setRootIndex(QModelIndex());
		updWindow();
	}
	actChangedToActive();
}

void Rightdemo_fs_Pannel::actChangedToActive()
{
	setActive();
	emit matchRightAsActive();
}

QVector< FileDescr > Rightdemo_fs_Pannel::getSelectedFiles()
{
	QVector< FileDescr > selectedRes;
	QModelIndexList selectList = listOfFiles->selectionModel()->selectedRows();

	if (selectList.isEmpty())
	{
		selectedRes.append(demo_fs_FileSystem->getFileDescript(listOfFiles->currentIndex()));
		return selectedRes;
	}

	for (const QModelIndex& ind : selectList)
	{
		selectedRes.append(demo_fs_FileSystem->getFileDescript(ind));
	}
	return selectedRes;
}

QByteArray Rightdemo_fs_Pannel::getFileContent(const FileDescr& fd)
{
	return demo_fs_FileSystem->getFileContent(fd);
}

void Rightdemo_fs_Pannel::actSortingByHeader(int headerInd)
{
	pannelsBase::sortByHeader(demo_fs_FileSystem, listOfFiles, lastHeaderFilterInd, sortingWay, headerInd);
	listOfFiles->reset();
	listOfFiles->setRootIndex(QModelIndex());
	actChangedToActive();
}

void Rightdemo_fs_Pannel::updWindow()
{
	listOfFiles->setRootIndex(QModelIndex());
	curLocationEdit->setText(getCurDir());
	backButton->setEnabled(pathHistory.backStack.size() > 1);
	forwardButton->setEnabled(!pathHistory.forwardStack.isEmpty());

	listOfFiles->reset();
}

void Rightdemo_fs_Pannel::setActive()
{
	listOfFiles->setFocus();
}

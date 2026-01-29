#include "leftpannel.h"

#include "utilspannels.h"

#include <QButtonGroup>
#include <QDebug>

LeftPannel::LeftPannel(QWidget* parent) : QWidget(parent)
{
	fileSystem = new LeftFileSystemModel(this);
	fileSystem->setRootPath(QDir::rootPath());
	;
	fileSystem->setHeaderData(3, Qt::Horizontal, "Time");

	QHeaderView* header;
	QHBoxLayout* pathSetting;
	QVBoxLayout* mainContainer;

	pannelsBase::initPannel(this, &listOfFiles, &curLocationEdit, &backButton, &forwardButton, &header, &pathSetting, &mainContainer);

	listOfFiles->setModel(fileSystem);

#ifdef Q_OS_WIN
	listOfFiles->setRootIndex(fileSystem->index(""));
	curPath = "";
#else
	listOfFiles->setRootIndex(fileSystem->index(QDir::homePath()));
	curPath = QDir::homePath();
#endif
	listOfFiles->sortByColumn(0, Qt::AscendingOrder);
	listOfFiles->setFocusPolicy(Qt::StrongFocus);
	listOfFiles->setColumnHidden(2, true);

	header->setSectionResizeMode(0, QHeaderView::Stretch);
	header->setSectionResizeMode(1, QHeaderView::ResizeToContents);
	header->setSectionResizeMode(2, QHeaderView::ResizeToContents);
	header->setSectionResizeMode(3, QHeaderView::ResizeToContents);

	curLocationEdit->setText(curPath);

	connect(listOfFiles, &QTreeView::doubleClicked, this, &LeftPannel::actItemSelected);
	connect(header, &QHeaderView::sectionClicked, [this](int col) { fileSystem->sort(col, Qt::AscendingOrder); });
	connect(header, &QHeaderView::sectionClicked, this, &LeftPannel::actSortingByHeader);

	connect(backButton, &QPushButton::clicked, this, &LeftPannel::returnOneBack);
	connect(forwardButton, &QPushButton::clicked, this, &LeftPannel::goOneForward);

	setLayout(mainContainer);
}

QString LeftPannel::getCurDir()
{
	return curPath;
}

QString LeftPannel::getFilePath()
{
	QModelIndex ind = listOfFiles->currentIndex();
	if (ind.isValid())
	{
		return fileSystem->filePath(ind);
	}
	return QString();
}

void LeftPannel::returnOneBack()
{
#ifdef Q_OS_WIN
	QModelIndex driveInd = fileSystem->index(QString());

	if (listOfFiles->rootIndex() == driveInd)
		return;

	if (pathHistory.backStack.isEmpty())
	{
		if (!listOfFiles->rootIndex().parent().isValid())
		{
			listOfFiles->setRootIndex(driveInd);
			curPath.clear();
		}

		curLocationEdit->setText(curPath);
		backButton->setEnabled(false);
		forwardButton->setEnabled(!pathHistory.forwardStack.isEmpty());
		actChangedToActive();
		return;
	}
#endif

	if (pathHistory.backStack.isEmpty())
		return;
	moveByOne(pathHistory.backStack, pathHistory.forwardStack);
}

void LeftPannel::goOneForward()
{
	if (pathHistory.forwardStack.isEmpty())
		return;
	moveByOne(pathHistory.forwardStack, pathHistory.backStack);
}

void LeftPannel::moveByOne(QStack< QString >& fromStack, QStack< QString >& toStack)
{
	toStack.push(curPath);
	curPath = fromStack.pop();
	updWindow();
	actChangedToActive();
}

void LeftPannel::actItemSelected(const QModelIndex& ind)
{
	if (!fileSystem->isDir(ind))
		return;
	setDir(fileSystem->filePath(ind));
	actChangedToActive();
}

void LeftPannel::actChangedToActive()
{
	setActive();
	emit matchLeftAsActive();
}

void LeftPannel::setActive()
{
	listOfFiles->setFocus();
}

void LeftPannel::setDir(const QString& setPath)
{
	QFileInfo check(setPath);
	if (!check.exists() || !check.isDir() || setPath == curPath)
		return;

	pathHistory.backStack.push(curPath);
	pathHistory.forwardStack.clear();
	curPath = setPath;
	updWindow();
}

void LeftPannel::actSortingByHeader(int headerInd)
{
	pannelsBase::sortByHeader(fileSystem, listOfFiles, lastHeaderFilterInd, sortingWay, headerInd);
	listOfFiles->reset();
	listOfFiles->setRootIndex(QModelIndex());
	actChangedToActive();
}

void LeftPannel::updWindow()
{
	QModelIndex ind = fileSystem->index(curPath);
	if (!ind.isValid())
		return;

	listOfFiles->setRootIndex(ind);
	curLocationEdit->setText(curPath);
	backButton->setEnabled(!pathHistory.backStack.isEmpty());
	forwardButton->setEnabled(!pathHistory.forwardStack.isEmpty());
}

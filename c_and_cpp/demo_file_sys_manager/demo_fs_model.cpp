#include "demo_fs_model.h"

#include <QBrush>
#include <QSet>

demo_fs_model::demo_fs_model(QObject *parent) : QAbstractItemModel(parent), curCluster(0)
{
	;
}

QModelIndex demo_fs_model::index(int row, int col, const QModelIndex &parent) const
{
	if (!hasIndex(row, col, parent))
		return QModelIndex();
	return createIndex(row, col);
}

int demo_fs_model::rowCount(const QModelIndex &parent) const
{
	if (parent.isValid())
		return 0;
	return entries.size();
}

int demo_fs_model::columnCount(const QModelIndex &parent) const
{
	return 3;
}

QVariant demo_fs_model::data(const QModelIndex &ind, int role) const
{
	if (!ind.isValid() || ind.row() >= entries.size())
		return QVariant();

	const FileDescr &entry = entries[ind.row()];
	if (role == Qt::DisplayRole)
	{
		switch (ind.column())
		{
		case 0:
			return entry.isActiveItem ? entry.name : "[DELETED]_" + entry.name;
		case 1:
			return (entry.type == "Dir") ? "Dir" : QString::number(entry.size);
		case 2:
			return entry.creationTime.toString("yyyy-MM-dd hh:mm:ss");
		}
	}

	if (!entry.isActiveItem && role == Qt::ForegroundRole)
		return QBrush(QColor(90, 130, 100, 140));

	return QVariant();
}

QVariant demo_fs_model::headerData(int ind, Qt::Orientation orientation, int role) const
{
	if (orientation != Qt::Horizontal || role != Qt::DisplayRole || ind < 0 || ind >= headNames.size())
		return QVariant();
	return headNames[ind];
}

struct FileDescriptComparator
{
	int col;
	Qt::SortOrder order;
	FileDescriptComparator(int c, Qt::SortOrder o) : col(c), order(o) {}

	bool operator()(const FileDescr &a, const FileDescr &b) const
	{
		bool asc = order == Qt::AscendingOrder;
		if (a.type == "Dir" && b.type != "Dir")
			return asc;
		if (a.type != "Dir" && b.type == "Dir")
			return !asc;
		if (col == 0)
			return asc ? a.name < b.name : b.name < a.name;
		if (col == 1)
			return asc ? a.size < b.size : b.size < a.size;
		if (col == 2)
			return asc ? a.creationTime < b.creationTime : b.creationTime < a.creationTime;
		return false;
	}
};

void demo_fs_model::sort(int col, Qt::SortOrder order)
{
	beginResetModel();
	std::sort(entries.begin(), entries.end(), FileDescriptComparator(col, order));
	endResetModel();
}

bool demo_fs_model::getdemo_fs_(const QString &path)
{
	if (!QFile::exists(path))
		return false;

	beginResetModel();
	entries.clear();
	if (path.isNull() || !demo_fs_Parser.unpackF(path))
	{
		endResetModel();
		return false;
	}
	curCluster = 0;
	updEntriesInfo();
	endResetModel();
	return true;
}

void demo_fs_model::setCurCluster(quint32 cluster)
{
	beginResetModel();
	curCluster = cluster;
	updEntriesInfo();
	endResetModel();
}

quint32 demo_fs_model::getCurCluster() const
{
	return curCluster;
}

FileDescr demo_fs_model::getFileDescript(const QModelIndex &index) const
{
	if (!index.isValid() || index.row() >= entries.size())
		return FileDescr();
	return entries[index.row()];
}

void demo_fs_model::updEntriesInfo()
{
	beginResetModel();
	entries = demo_fs_Parser.dirList(curCluster);
	endResetModel();
}

QString demo_fs_model::formPath()
{
	if (curCluster == demo_fs_Parser.getRootCluster())
		return "/";

	QStringList pathFormer;
	quint32 cur = curCluster;
	QSet< quint32 > visited;
	while ((demo_fs_Parser.isdemo_fs_32Type() ? cur != demo_fs_Parser.getRootCluster() : cur != 0) && !visited.contains(cur))
	{
		visited.insert(cur);
		quint32 parentCluster = demo_fs_Parser.getParentCluster(cur);
		QVector< FileDescr > curUpperDir = demo_fs_Parser.dirList(parentCluster);

		bool found = false;
		for (const FileDescr &fd : curUpperDir)
		{
			if (fd.type != "Dir" || fd.clusterStart != cur || fd.name == "..")
				continue;

			pathFormer.prepend(fd.name);
			found = true;
			break;
		}
		if (!found)
			break;

		cur = parentCluster;
	}
	QString path = pathFormer.isEmpty() ? "/" : "/" + pathFormer.join("/");
	return path;
}

QModelIndex demo_fs_model::parent(const QModelIndex &index) const
{
	return QModelIndex();
}

bool demo_fs_model::isdemo_fs_32TypeModel() const
{
	return demo_fs_Parser.isdemo_fs_32Type();
}
quint32 demo_fs_model::getRootClusterModel() const
{
	return demo_fs_Parser.getRootCluster();
}

QByteArray demo_fs_model::getFileContent(const FileDescr &fd)
{
	return demo_fs_Parser.readFileContent(fd);
}
